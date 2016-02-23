package system

import(
    "os"
    "io"
    "encoding/csv"
    "log"
    "strconv"
)

type PasswdEntry struct {
    Name string
    Uid int
    Gid int
    Shell string
    Fullname string
    Home string
}

type OSUser struct {
    ShadowEntry []string
    PasswdEntry []string
}

type OSGroup struct {
    GroupEntry []string
}

func LoadUsers(root string) []*OSUser {
  var users []*OSUser

  shadowf, serr := os.Open(path.Join(root, "etc", "shadow"))
  passwdf, perr := os.Open(path.Join(root, "etc", "passwd"))

  var pwentries, shadowentries [][]string

  if serr != nil || perr != nil {
    return users
  }

  shadowreader := csv.NewReader(shadowf)
  shadowreader.Comma = ':'

  for {
    record, err := shadowreader.Read()
    if err != nil {
      if err == io.EOF {
        break
      }
      log.Fatal(err)
    }

    append(shadowentries, record)
  }

  passwdreader := csv.NewReader(passwdf)
  passwdreader.Comma = ':'

  for {
    record, err := passwdreader.Read()
    if err != nil {
      if err == io.EOF {
        break
      }
      log.Fatal(err)
    }

    for shadowentry := range shadowentries {
      if shadowentry[0] == record[0] {
        append(users, OSUser{shadowentry, entry})
        break
      }
    }
  }
}


func ReadPasswdEntry(r []string) (p *PasswdEntry) {
    p = new(PasswdEntry)

    p.Name = r[0]
    p.Uid, _ = strconv.Atoi(r[2])
    p.Gid, _ = strconv.Atoi(r[3])
    p.Fullname = r[4]
    p.Home = r[5]
    p.Shell = r[6]

    return p
}

func ReadEtcPasswd(path string) []*PasswdEntry {
    var entries []*PasswdEntry

    f, err := os.Open(path)

    if err != nil {
        return entries
    }

    csvreader := csv.NewReader(f)
    csvreader.Comma = ':'

    for {
        record, err := csvreader.Read()
        if err != nil {
            if err == io.EOF {
                break
            }
            log.Fatal(err)
        }

        entries = append(entries, ReadPasswdEntry(record))
    }

    return entries
}

func (c Container) AddHostUser(u system.SystemUser) {
    etcShadowPath := path.Join(c.Location, c.Name, "etc", "shadow")
    etcPasswdPath := path.Join(c.Location, c.Name, "etc", "passwd")
    //etcGroupsPath := path.Join(c.Location, c.Name, "etc", "groups")

    uids := strconv.Itoa(u.Uid)
    gids := strconv.Itoa(u.Gid)

    shadow, _ := os.OpenFile(etcShadowPath, os.O_APPEND|os.O_WRONLY, 0600)
    shadow.WriteString(u.Name + ":" + u.Pwhash + ":" + "14675" + ":" + "99999" + ":" + "99999" + ":" + "7" + ":::")

    passwd, _ := os.OpenFile(etcPasswdPath, os.O_APPEND|os.O_WRONLY, 0600)
    passwd.WriteString(u.Name + ":x:" + uids + ":" + gids + "::" + "/home/" + u.Name + ":" + u.Shell)
}
