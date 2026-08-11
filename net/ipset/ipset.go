package ipset

import (
	"context"
	"fmt"
	"log"
	"strings"

	knftables "sigs.k8s.io/knftables"
)

type Name string
type Type string
type UID string

const (
	ListSet = Type("list:set")
	HashIP  = Type("hash:ip")
	HashNet = Type("hash:net")

	nftTableName = "weave"
)

type Interface interface {
	Create(ipsetName Name, ipsetType Type) error
	AddEntry(user UID, ipsetName Name, entry string, comment string) error
	DelEntry(user UID, ipsetName Name, entry string) error
	EntryExists(user UID, ipsetName Name, entry string) bool
	Exists(ipsetName Name) (bool, error)
	Flush(ipsetName Name) error
	Destroy(ipsetName Name) error
	List(prefix string) ([]Name, error)
	FlushAll() error
	DestroyAll() error
}

type entryKey struct {
	ipsetName Name
	entry     string
}

type ipset struct {
	*log.Logger
	maxListSize int
	users       map[entryKey]map[UID]struct{}
	nft         knftables.Interface
	initErr     error
}

func New(logger *log.Logger, maxListSize int) Interface {
	nft, err := knftables.New(knftables.InetFamily, nftTableName)
	if err != nil {
		return &ipset{Logger: logger, maxListSize: maxListSize, users: make(map[entryKey]map[UID]struct{}), initErr: err}
	}
	return newWithInterface(logger, maxListSize, nft)
}

func newWithInterface(logger *log.Logger, maxListSize int, nft knftables.Interface) Interface {
	return &ipset{
		Logger:      logger,
		maxListSize: maxListSize,
		users:       make(map[entryKey]map[UID]struct{}),
		nft:         nft,
	}
}

func nftSetName(name Name) string {
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

func (i *ipset) ready() error {
	if i.initErr != nil {
		return fmt.Errorf("initialize nftables sets: %w", i.initErr)
	}
	tx := i.nft.NewTransaction()
	tx.Add(&knftables.Table{})
	if err := i.nft.Run(context.Background(), tx); err != nil && !knftables.IsAlreadyExists(err) {
		return fmt.Errorf("create nftables table: %w", err)
	}
	return nil
}

func (i *ipset) Create(ipsetName Name, ipsetType Type) error {
	if ipsetType == ListSet {
		return fmt.Errorf("nftables does not support nested %q sets; namespace selectors must contain pod IPs directly", ListSet)
	}
	if ipsetType != HashIP && ipsetType != HashNet {
		return fmt.Errorf("unsupported nftables set type %q", ipsetType)
	}
	if err := i.ready(); err != nil {
		return err
	}
	exists, err := i.Exists(ipsetName)
	if err != nil || exists {
		return err
	}
	set := &knftables.Set{Name: nftSetName(ipsetName), Type: "ipv4_addr"}
	if ipsetType == HashNet {
		set.Flags = []knftables.SetFlag{knftables.IntervalFlag}
		autoMerge := true
		set.AutoMerge = &autoMerge
	}
	if i.maxListSize > 0 {
		size := uint64(i.maxListSize)
		set.Size = &size
	}
	tx := i.nft.NewTransaction()
	tx.Create(set)
	if err := i.nft.Run(context.Background(), tx); err != nil && !knftables.IsAlreadyExists(err) {
		return fmt.Errorf("create nftables set %s: %w", ipsetName, err)
	}
	return nil
}

func (i *ipset) AddEntry(user UID, ipsetName Name, entry string, comment string) error {
	i.Logger.Printf("adding entry %s to %s for %s", entry, ipsetName, user)
	if !i.addUser(user, ipsetName, entry) {
		return nil
	}
	element := &knftables.Element{Set: nftSetName(ipsetName), Key: []string{entry}}
	if comment != "" {
		element.Comment = &comment
	}
	tx := i.nft.NewTransaction()
	tx.Add(element)
	if err := i.nft.Run(context.Background(), tx); err != nil {
		i.delUser(user, ipsetName, entry)
		return fmt.Errorf("add %s to nftables set %s: %w", entry, ipsetName, err)
	}
	return nil
}

func (i *ipset) DelEntry(user UID, ipsetName Name, entry string) error {
	i.Logger.Printf("deleting entry %s from %s for %s", entry, ipsetName, user)
	if !i.delUser(user, ipsetName, entry) {
		return nil
	}
	tx := i.nft.NewTransaction()
	tx.Delete(&knftables.Element{Set: nftSetName(ipsetName), Key: []string{entry}})
	if err := i.nft.Run(context.Background(), tx); err != nil && !knftables.IsNotFound(err) {
		return fmt.Errorf("delete %s from nftables set %s: %w", entry, ipsetName, err)
	}
	return nil
}

func (i *ipset) EntryExists(user UID, ipsetName Name, entry string) bool {
	return i.existUser(user, ipsetName, entry)
}

func (i *ipset) Exists(name Name) (bool, error) {
	sets, err := i.nft.List(context.Background(), "set")
	if err != nil {
		if knftables.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	wanted := nftSetName(name)
	for _, set := range sets {
		if set == wanted {
			return true, nil
		}
	}
	return false, nil
}

func (i *ipset) Flush(ipsetName Name) error {
	i.removeSetFromUsers(ipsetName)
	tx := i.nft.NewTransaction()
	tx.Flush(&knftables.Set{Name: nftSetName(ipsetName)})
	if err := i.nft.Run(context.Background(), tx); err != nil && !knftables.IsNotFound(err) {
		return err
	}
	return nil
}

func (i *ipset) Destroy(ipsetName Name) error {
	i.removeSetFromUsers(ipsetName)
	tx := i.nft.NewTransaction()
	tx.Destroy(&knftables.Set{Name: nftSetName(ipsetName)})
	if err := i.nft.Run(context.Background(), tx); err != nil && !knftables.IsNotFound(err) {
		return err
	}
	return nil
}

func (i *ipset) List(prefix string) ([]Name, error) {
	sets, err := i.nft.List(context.Background(), "set")
	if err != nil {
		if knftables.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	normalizedPrefix := nftSetName(Name(prefix))
	selected := make([]Name, 0)
	for _, set := range sets {
		if strings.HasPrefix(set, normalizedPrefix) {
			selected = append(selected, Name(set))
		}
	}
	return selected, nil
}

func (i *ipset) FlushAll() error {
	sets, err := i.List("")
	if err != nil {
		return err
	}
	for _, set := range sets {
		if err := i.Flush(set); err != nil {
			return err
		}
	}
	return nil
}

func (i *ipset) DestroyAll() error {
	sets, err := i.List("")
	if err != nil {
		return err
	}
	for _, set := range sets {
		if err := i.Destroy(set); err != nil {
			return err
		}
	}
	return nil
}

func (i *ipset) addUser(user UID, ipsetName Name, entry string) bool {
	key := entryKey{ipsetName, entry}
	if i.users[key] == nil {
		i.users[key] = make(map[UID]struct{})
	}
	add := len(i.users[key]) == 0
	i.users[key][user] = struct{}{}
	return add
}

func (i *ipset) delUser(user UID, ipsetName Name, entry string) bool {
	key := entryKey{ipsetName, entry}
	owners := i.users[key]
	if _, found := owners[user]; !found {
		return false
	}
	delete(owners, user)
	if len(owners) != 0 {
		return false
	}
	delete(i.users, key)
	return true
}

func (i *ipset) existUser(user UID, ipsetName Name, entry string) bool {
	_, ok := i.users[entryKey{ipsetName, entry}][user]
	return ok
}

func (i *ipset) removeSetFromUsers(ipsetName Name) {
	for key := range i.users {
		if key.ipsetName == ipsetName {
			delete(i.users, key)
		}
	}
}
