package OUI

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"fmt"
	"errors"
)

const (
	macLength = 6
)

/**
 * 错误信息
 */
var ErrInvalidMACAddress = errors.New("invalid MAC address")

/**
 * oui []byte(mac)
 * mac mac地址
 * Organization 组织机构
 * Organization2 内容组织机构在收集范围中(InsideOrganization)的对照名
 */
type AddressBlock struct {
	Oui          []byte
	OuiMac       string
	Organization string
	Organization2 string
}

// 地图指针
type mapAddress struct {
	Block *AddressBlock
}

/**
 * 组织机构统一大写
 * InsideOrganization: 收集在内的数据要在这个范围中. *%s*
 * OutsideOrganization: 排除在外的数据
*/
type OuiDb struct {
	ioz []string
	ooz []string

	Blocks []AddressBlock
	Maps   map[string]mapAddress
}

var db = &OuiDb{
	nil,
	nil,
	nil,
	nil,
}

/**
 * 创建 oui db
 */
func (m *OuiDb) Open(file string) *OuiDb {
	if err := db.load(file); err != nil {
		return nil
	}
	return db
}

/**
 * 设置收录数据
*/
func (m *OuiDb) SetInsideOrganization(s []string)  {
	db.ioz = s
}

/**
 * 排除数据
*/
func (m *OuiDb) SetOutsideOrganization(s []string)  {
	db.ooz = s
}

/**
 * 加载 OUI 文件
 * @param file path
 */
func (m *OuiDb) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	var organization2 string
	var regS string

	fieldsRe := regexp.MustCompile(`^(\S+)\s+\(base 16\)\s+(.+)?`)

	scanner := bufio.NewScanner(file)
	F:
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || text[0] == '#' || text[0] == '\t' || strings.Index(text, "(base 16)") < 0 {
			continue
		}

		fields := fieldsRe.FindAllStringSubmatch(text, -1)

		addr := fields[0][1]
		organization := strings.ToUpper(fields[0][2])

		if m.ooz == nil &&
			m.ioz == nil {
			goto add
		}
		organization2 = ""

		if m.ooz != nil {
			for i := len(m.ooz) - 1; i >= 0; i-- {
				o := m.ooz[i]
				regS = fmt.Sprintf(`\b%s\b`, o)
				if f, _ := regexp.MatchString(regS, organization); f {
					continue F
				}
			}
			goto add
		}
		if m.ioz != nil {
			for i := len(m.ioz) - 1; i >= 0; i-- {
				o := m.ioz[i]
				regS = fmt.Sprintf(`\b%s\b`, o)
				if f, _ := regexp.MatchString(regS, organization); f && organization != ""  {
					organization2 = o
					goto add
				}
			}
		}
		continue
	add:
		m.Blocks = append(m.Blocks, AddressBlock{
			Oui: []byte(addr),
			OuiMac: addr,
			Organization: organization,
			Organization2: organization2,
		})
	}

	m.createSmartMap()

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

/**
 * 创建指针地图指引
 */
func (m *OuiDb) createSmartMap() {
	m.Maps = make(map[string]mapAddress)
	for i := len(m.Blocks) - 1; i >= 0; i-- {
		obj := m.Blocks[i]
		m.Maps[obj.OuiMac] = mapAddress{
			&m.Blocks[i],
		}
	}
}

/**
 * 获取信息
 */
func (m *OuiDb) VendorLookup(mac string) (ab *AddressBlock, err error) {
	mac = strings.Replace(strings.Replace(mac, "-", "", -1),":","",-1)

	if len(mac) < macLength {
		goto errorFn
	}
	if k, ok := m.Maps[mac[:macLength]]; ok {
		return k.Block, nil
	} else {
		goto errorFn
	}

errorFn:
	return nil, ErrInvalidMACAddress
}
