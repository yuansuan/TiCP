package migration

import (
	"fmt"
	"io/fs"
	"regexp"
	"strconv"
	"strings"
)

type DatabaseType string

const (
	Mysql DatabaseType = "mysql"
)

func (dt DatabaseType) String() string {
	return string(dt)
}

type migrationType string

const (
	up   migrationType = "up"
	down migrationType = "down"
)

func parseToMigrationType(s string) (migrationType, error) {
	if s == string(up) || s == string(down) {
		return migrationType(s), nil
	}

	return "", fmt.Errorf("unknown migration type, %s", s)
}

type sourceNode struct {
	version   int
	upgrade   *sourceFile
	downgrade *sourceFile

	prev *sourceNode
	next *sourceNode
}

func (sn *sourceNode) clone() *sourceNode {
	n := &sourceNode{
		version: sn.version,
	}

	if sn.upgrade != nil {
		n.upgrade = sn.upgrade.clone()
	}

	if sn.downgrade != nil {
		n.downgrade = sn.downgrade.clone()
	}

	return n
}

func newSourceNode(version int, mType migrationType, sqlScript []byte, filename string) *sourceNode {
	sn := &sourceNode{
		version: version,
	}

	sf := &sourceFile{
		sqlScript: sqlScript,
		filename:  filename,
	}

	switch mType {
	case up:
		sn.upgrade = sf
	case down:
		sn.downgrade = sf
	}

	return sn
}

type sourceFile struct {
	sqlScript []byte
	filename  string
}

func (sf *sourceFile) clone() *sourceFile {
	return &sourceFile{
		sqlScript: sf.sqlScript,
		filename:  sf.filename,
	}
}

type source struct {
	root *sourceNode
}

func newSource(migrationSourceFS fs.FS, dbType DatabaseType) (*source, error) {
	s := &source{
		root: &sourceNode{
			version: 0,
		},
	}

	dirEntries, err := fs.ReadDir(migrationSourceFS, dbType.String())
	if err != nil {
		return nil, fmt.Errorf("migration source fs read dir failed, %w", err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		// check filename match like 1.init.up.sql
		version, mType, err := parseSourceFilename(dirEntry.Name())
		if err != nil {
			return nil, fmt.Errorf("parse source file name failed, %w", err)
		}

		// error occur in windows when using filepath.Join() so chose to concat filePath by string
		filePath := dbType.String() + "/" + dirEntry.Name()
		sqlScript, err := fs.ReadFile(migrationSourceFS, filePath)
		if err != nil {
			return nil, fmt.Errorf("read file %s failed, %w", filePath, err)
		}

		// check if version migration exist, cannot multiple append
		exist := s.findAndInsert(version, mType, sqlScript, dirEntry.Name())
		if exist {
			return nil, fmt.Errorf("duplicated version script file in migration type [%s] about version %d", mType, version)
		}
	}

	// check if source link list valid
	if err = s.checkSourceLinkListValid(); err != nil {
		return nil, fmt.Errorf("migration source files invalid, %w", err)
	}

	return s, nil
}

func (s *source) findAndInsert(version int, mType migrationType, sqlScript []byte, filename string) (exist bool) {
	curr := s.root

	var tail *sourceNode
	for curr != nil {
		if curr.version == version {
			switch mType {
			case up:
				if curr.upgrade != nil {
					return true
				}

				curr.upgrade = &sourceFile{
					sqlScript: sqlScript,
					filename:  filename,
				}
			case down:
				if curr.downgrade != nil {
					return true
				}

				curr.downgrade = &sourceFile{
					sqlScript: sqlScript,
					filename:  filename,
				}
			}

			return false
		} else if curr.version > version { // insert to before
			newNode := newSourceNode(version, mType, sqlScript, filename)
			if curr.prev != nil {
				curr.prev.next = newNode
				newNode.prev = curr.prev
			}

			newNode.next = curr
			curr.prev = newNode

			return false
		}

		tail = curr
		curr = curr.next
	}

	newNode := newSourceNode(version, mType, sqlScript, filename)
	tail.next = newNode
	newNode.prev = tail

	return false
}

const migrationSourceFileRegPattern = `^\d+\.[a-zA-Z0-9_-]+\.(up|down)\.sql$`

func parseSourceFilename(filename string) (int, migrationType, error) {
	// must be like 1.init.up.sql / 1.init.down.sql
	match, err := regexp.MatchString(migrationSourceFileRegPattern, filename)
	if err != nil {
		return 0, "", fmt.Errorf("regex error, %w", err)
	}

	if match {
		fields := strings.Split(filename, ".")
		if len(fields) != 4 {
			return 0, "", fmt.Errorf("split with . not return 4 fields, filename [%s]", filename)
		}

		version, err := strconv.Atoi(fields[0])
		if err != nil {
			return 0, "", fmt.Errorf("parse first fields to int failed, filename [%s], %w", filename, err)
		}

		mType, err := parseToMigrationType(fields[2])
		if err != nil {
			return 0, "", fmt.Errorf("parse %s to migration type failed, %w", fields[2], err)
		}

		return version, mType, nil
	} else {
		return 0, "", fmt.Errorf("%s not match %s", filename, migrationSourceFileRegPattern)
	}
}

func (s *source) checkSourceLinkListValid() error {
	if s.root == nil {
		return fmt.Errorf("root node is nil")
	}

	curr := s.root.next
	currentVersionShouldBe := 0
	for curr != nil {
		currentVersionShouldBe++
		if curr.version != currentVersionShouldBe {
			return fmt.Errorf("lack of version script file in version %d", currentVersionShouldBe)
		}

		if curr.upgrade == nil {
			return fmt.Errorf("lack of upgrade script in version %d", curr.version)
		}

		if curr.downgrade == nil {
			return fmt.Errorf("lack of downgrade script in version %d", curr.version)
		}

		curr = curr.next
	}

	return nil
}

func (s *source) getLatestVersion() int {
	curr := s.root

	for curr.next != nil {
		curr = curr.next
	}

	if curr.version < 0 {
		return 0
	}

	return curr.version
}

func (s *source) subLinkList(from, to int) (*sourceNode, *sourceNode, error) {
	if from >= to {
		return nil, nil, fmt.Errorf("sub link list from [%d] should less than to [%d]", from, to)
	}

	var head *sourceNode
	var tail *sourceNode
	curr := s.root
	for curr != nil {
		if curr.version == from {
			head = curr.clone()
			tail = head
		}

		if from < curr.version && curr.version <= to {
			nodeClone := curr.clone()
			tail.next = nodeClone
			nodeClone.prev = tail
			tail = tail.next
		}

		// at the tail, check if tail version less than from/to
		if curr.next == nil {
			if curr.version < from {
				return nil, nil, fmt.Errorf("the max version is %d while want to sub from %d", curr.version, from)
			}

			if curr.version < to {
				return nil, nil, fmt.Errorf("the max version is %d while want to sub to %d", curr.version, to)
			}
		}

		curr = curr.next
	}

	return head, tail, nil
}
