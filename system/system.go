package system

import(
    "os"
    "io"
    "encoding/csv"
    "log"
    "regexp"
    "bufio"
    "strings"
    "strconv"
    "fmt"
    "path"
)

type OSUser struct {
    shadowEntry []string
    passwdEntry []string
}

func (o OSUser) Username() string {
  return o.shadowEntry[0]
}

func (o OSUser) Uid() int {
  uid, _ := strconv.Atoi(o.passwdEntry[2])
  return uid
}

func (o OSUser) Gid() int {
  gid, _ := strconv.Atoi(o.passwdEntry[3])
  return gid
}

func (o OSUser) Home() string {
  return o.passwdEntry[5]
}

func (o OSUser) Shell() string {
  return o.passwdEntry[6]
}

func (o OSUser) SetShell(sh string) {
  o.passwdEntry[6] = sh
}

type OSUsers []*OSUser

func (o OSUsers) Find(username string) *OSUser {
  var user *OSUser

  for _, u := range o {
    if u.Username() == username {
      user = u
      break
    }
  }

  return user
}

func (o OSUser) BuildEtcShadowEntry() string {
  return strings.Join(o.shadowEntry, ":")
}

func (o OSUser) BuildEtcPasswdEntry() string {
  return strings.Join(o.passwdEntry, ":")
}

func (o OSUser) CreateUser(root string) error {
  shell := o.Shell()

  _, err := os.Stat(path.Join(root, shell))
  if os.IsNotExist(err) {
    o.SetShell("/bin/bash")
  }

  err = o.UpdateEntry(root)

  if err != nil {
    return err
  }

  homedir := o.Home()
  err = os.MkdirAll(path.Join(root, homedir), 0700)
  err = os.Chown(path.Join(root, homedir), o.Uid(), o.Gid())

  return err
}

func (o OSUser) UpdateEntry(root string) error {
  shadowpath := path.Join(root, "etc", "shadow")
  passwdpath := path.Join(root, "etc", "passwd")

  fmt.Println(shadowpath)

  shadowf, serr := os.OpenFile(shadowpath, os.O_RDWR|os.O_APPEND, 0600)
  passwdf, perr := os.OpenFile(passwdpath, os.O_RDWR|os.O_APPEND, 0600)
  defer shadowf.Close()
  defer passwdf.Close()

  if serr != nil {
    return serr
  }

  if perr != nil {
    return perr
  }

  reg := regexp.MustCompile("^" + o.Username() + ":.*")
  found := reg.MatchReader(bufio.NewReader(shadowf)) || reg.MatchReader(bufio.NewReader(passwdf))

  if found {
    return fmt.Errorf("User %s already exists in container at %s", o.Username(), root)
  }

  _, serr = shadowf.WriteString(o.BuildEtcShadowEntry())
  _, perr = passwdf.WriteString(o.BuildEtcPasswdEntry())

  if serr != nil || perr != nil {
    return fmt.Errorf("Error writing user: %s %s", serr, perr)
  }

  return nil

}

func LoadUsers(root string) OSUsers {
  var users OSUsers

  shadowf, serr := os.Open(path.Join(root, "etc", "shadow"))
  passwdf, perr := os.Open(path.Join(root, "etc", "passwd"))
  defer shadowf.Close()
  defer passwdf.Close()

  var shadowentries [][]string

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

    shadowentries = append(shadowentries, record)
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

    for _, shadowentry := range shadowentries {
      if shadowentry[0] == record[0] {
        found_user := new(OSUser)
        found_user.shadowEntry = shadowentry
        found_user.passwdEntry = record
        users = append(users, found_user)
        break
      }
    }
  }

  return users
}
