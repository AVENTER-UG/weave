package nftables

import (
	"fmt"
	"strconv"
	"strings"
)

func sanitizeName(name string) string {
	var sanitized strings.Builder
	sanitized.Grow(len(name))
	for _, character := range name {
		switch {
		case character >= 'a' && character <= 'z',
			character >= 'A' && character <= 'Z',
			character >= '0' && character <= '9',
			character == '_':
			sanitized.WriteRune(character)
		default:
			sanitized.WriteByte('_')
		}
	}
	return sanitized.String()
}

func chainName(table, chain string) string {
	return sanitizeName(table + "_" + chain)
}

func setName(name string) string {
	return sanitizeName(name)
}

func parseMark(value string) (uint32, uint32, error) {
	parts := strings.SplitN(value, "/", 2)
	mark, err := strconv.ParseUint(parts[0], 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid mark %q: %w", value, err)
	}
	mask := uint64(0xffffffff)
	if len(parts) == 2 {
		mask, err = strconv.ParseUint(parts[1], 0, 32)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid mark mask %q: %w", value, err)
		}
	}
	return uint32(mark), uint32(mask), nil
}

func translateRule(table string, args []string) (string, error) {
	var (
		parts     []string
		negated   bool
		protocol  string
		target    string
		policyDir string
	)

	needValue := func(i int, option string) (string, error) {
		if i+1 >= len(args) {
			return "", fmt.Errorf("missing value for %s", option)
		}
		return args[i+1], nil
	}
	comparison := func() string {
		if negated {
			negated = false
			return " != "
		}
		return " "
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "!":
			negated = true
		case "-m":
			if _, err := needValue(i, arg); err != nil {
				return "", err
			}
			i++ // nft expressions do not need iptables match modules
		case "-i", "--in-interface":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			parts = append(parts, "iifname"+comparison()+strconv.Quote(value))
			i++
		case "-o", "--out-interface":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			parts = append(parts, "oifname"+comparison()+strconv.Quote(value))
			i++
		case "-s", "--source", "--src":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			parts = append(parts, "ip saddr"+comparison()+value)
			i++
		case "-d", "--destination", "--dst":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			parts = append(parts, "ip daddr"+comparison()+value)
			i++
		case "-p", "--protocol":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			protocol = strings.ToLower(value)
			parts = append(parts, "meta l4proto"+comparison()+protocol)
			i++
		case "--dport", "--destination-port":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			transport := protocol
			if transport != "tcp" && transport != "udp" && transport != "sctp" {
				transport = "th"
			}
			parts = append(parts, transport+" dport"+comparison()+value)
			i++
		case "--sport", "--source-port":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			transport := protocol
			if transport != "tcp" && transport != "udp" && transport != "sctp" {
				transport = "th"
			}
			parts = append(parts, transport+" sport"+comparison()+value)
			i++
		case "--state", "--ctstate":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			states := strings.Split(strings.ToLower(value), ",")
			expr := states[0]
			if len(states) > 1 {
				expr = "{ " + strings.Join(states, ", ") + " }"
			}
			parts = append(parts, "ct state"+comparison()+expr)
			i++
		case "--match-set":
			name, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			if i+2 >= len(args) {
				return "", fmt.Errorf("missing direction for %s", arg)
			}
			direction := args[i+2]
			if direction != "src" && direction != "dst" {
				return "", fmt.Errorf("unsupported set direction %q", direction)
			}
			addressKey := "daddr"
			if direction == "src" {
				addressKey = "saddr"
			}
			parts = append(parts, "ip "+addressKey+comparison()+"@"+setName(name))
			i += 2
		case "--mark":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			mark, mask, err := parseMark(value)
			if err != nil {
				return "", err
			}
			op := "=="
			if negated {
				op = "!="
				negated = false
			}
			parts = append(parts, fmt.Sprintf("meta mark & 0x%x %s 0x%x", mask, op, mark))
			i++
		case "--set-xmark":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			mark, mask, err := parseMark(value)
			if err != nil {
				return "", err
			}
			parts = append(parts, fmt.Sprintf("meta mark set meta mark & 0x%x | 0x%x", ^mask, mark&mask))
			i++
		case "--espspi":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			parts = append(parts, "esp spi"+comparison()+value)
			i++
		case "--dir":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			policyDir = value
			i++
		case "--pol":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			if value != "none" || (policyDir != "in" && policyDir != "out") {
				return "", fmt.Errorf("unsupported IPsec policy %s/%s", policyDir, value)
			}
			if policyDir == "out" {
				parts = append(parts, "rt ipsec missing")
			} else {
				parts = append(parts, "meta ipsec missing")
			}
			i++
		case "--physdev-is-bridged":
			// A following --physdev-in/out match proves this condition. The few
			// rules without an interface are translated when that case is needed.
		case "--physdev-in", "--physdev-out":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			key := "meta bri_iifname"
			if arg == "--physdev-out" {
				key = "meta bri_oifname"
			}
			parts = append(parts, key+comparison()+strconv.Quote(value))
			i++
		case "--src-type", "--dst-type":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			direction := "saddr"
			if arg == "--dst-type" {
				direction = "daddr"
			}
			parts = append(parts, "fib "+direction+" type"+comparison()+strings.ToLower(value))
			i++
		case "--comment":
			if _, err := needValue(i, arg); err != nil {
				return "", err
			}
			i++ // comments are represented by the backend's rule identity
		case "-j", "--jump":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			target = value
			i++
		case "--nflog-group":
			value, err := needValue(i, arg)
			if err != nil {
				return "", err
			}
			parts = append(parts, "log group "+value)
			i++
		default:
			if strings.HasPrefix(arg, "--physdev-in=") || strings.HasPrefix(arg, "--physdev-out=") {
				pieces := strings.SplitN(arg, "=", 2)
				key := "meta bri_iifname"
				if pieces[0] == "--physdev-out" {
					key = "meta bri_oifname"
				}
				parts = append(parts, key+comparison()+strconv.Quote(pieces[1]))
				continue
			}
			return "", fmt.Errorf("unsupported iptables argument %q in %v", arg, args)
		}
	}

	if negated {
		return "", fmt.Errorf("dangling negation in %v", args)
	}

	switch strings.ToUpper(target) {
	case "":
	case "ACCEPT":
		parts = append(parts, "accept")
	case "DROP":
		parts = append(parts, "drop")
	case "RETURN":
		parts = append(parts, "return")
	case "MASQUERADE":
		parts = append(parts, "masquerade")
	case "NFLOG", "MARK":
		// Their statements were emitted while parsing target options.
	default:
		parts = append(parts, "jump "+chainName(table, target))
	}

	return strings.Join(parts, " "), nil
}
