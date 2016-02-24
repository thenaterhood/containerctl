package system

import(
    "os"
    "io"
    "encoding/csv"
    "log"
    "regexp"
    "bufio"
    "strings"
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

  fmt.Println("Checking if user exists")

  reg := regexp.MustCompile("^" + o.Username() + ":.*")
  found := reg.MatchReader(bufio.NewReader(shadowf)) || reg.MatchReader(bufio.NewReader(passwdf))

  if found {
    return fmt.Errorf("User %s already exists in container at %s", o.Username(), root)
  }

  fmt.Println("Appending user")

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
