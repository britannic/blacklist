# edgeos

--
    import "github.com/britannic/blacklist/internal/edgeos"

Package edgeos provides methods and structures to retrieve, parse and render
EdgeOS configuration data and files.

## Usage

```go
const (

 // ExcDomns is a string labels for domain exclusions
 ExcDomns = "whitelisted-subdomains"
 // ExcHosts is a string labels for host exclusions
 ExcHosts = "whitelisted-servers"
 // ExcRoots is a string labels for preconfigured global domain exclusions
 ExcRoots = "global-whitelisted-domains"
 // PreDomns is a string label for preconfigured whitelisted domains
 PreDomns = "blacklisted-subdomains"
 // PreHosts is a string label for preconfigured blacklisted hosts
 PreHosts = "blacklisted-servers"
 // PreRoots is a string label for preconfigured global blacklisted hosts
 PreRoots = "global-blacklisted-domains"
 // False is a string constant
 False = "false"
 // True is a string constant
 True = "true"
)
```

#### func  ChkWeb

```go
func ChkWeb(site string, port int) bool
```

ChkWeb returns true if DNS is working

#### func  GetFile

```go
func GetFile(f string) (io.Reader, error)
```

GetFile reads a file and returns an io.Reader

#### func  Iter

```go
func Iter(i int) []struct{}
```

Iter iterates over ints - use it in for loops

#### func  NewWriter

```go
func NewWriter() io.Writer
```

NewWriter returns an io.Writer

#### type CFGcli

```go
type CFGcli struct {
 *Config
 Cfg string
}
```

CFGcli loads configurations using the EdgeOS CFGcli

#### type CFGstatic

```go
type CFGstatic struct {
 *Config
 Cfg string
}
```

CFGstatic loads static configurations for testing

#### type CFile

```go
type CFile struct {
 *Env
 Names []string
}
```

CFile holds an array of file names

#### func (*CFile) Remove

```go
func (c *CFile) Remove() error
```

Remove deletes a CFile array of file names

#### func (*CFile) String

```go
func (c *CFile) String() string
```

String implements string method

#### func (*CFile) Strings

```go
func (c *CFile) Strings() []string
```

Strings returns a sorted array of strings.

#### type ConfLoader

```go
type ConfLoader interface {
 // contains filtered or unexported methods
}
```

ConfLoader interface handles multiple configuration load methods

#### type Config

```go
type Config struct {
 *Env
}
```

Config is a struct of configuration fields

#### func  NewConfig

```go
func NewConfig(opts ...Option) *Config
```

NewConfig returns a new *Config initialized with the parameter options passed to
it

#### func (*Config) Blacklist

```go
func (c *Config) Blacklist(r ConfLoader) error
```

Blacklist extracts blacklist nodes from a EdgeOS/VyOS configuration structure

#### func (*Config) Get

```go
func (c *Config) Get(nx string) *Objects
```

Get returns an *Object for a given node

#### func (*Config) GetAll

```go
func (c *Config) GetAll(ltypes ...string) *Objects
```

GetAll returns a pointer to an Objects struct

#### func (*Config) GetTotalStats

```go
func (c *Config) GetTotalStats() (dropped, extracted, kept int32)
```

GetTotalStats displays aggregate statistics for processed sources

#### func (*Config) InSession

```go
func (c *Config) InSession() bool
```

InSession returns true if VyOS/EdgeOS configure is in session

#### func (*Config) NewContent

```go
func (c *Config) NewContent(iface IFace) (Contenter, error)
```

NewContent returns a Contenter interface of the requested IFace type

#### func (*Config) Nodes

```go
func (c *Config) Nodes() (n []string)
```

Nodes returns an array of configured nodes

#### func (*Config) ProcessContent

```go
func (c *Config) ProcessContent(cts ...Contenter) error
```

ProcessContent processes the Contents array

#### func (*Config) ReloadDNS

```go
func (c *Config) ReloadDNS() ([]byte, error)
```

ReloadDNS reloads the dnsmasq configuration

#### func (*Config) SetOpt

```go
func (c *Config) SetOpt(opts ...Option) Option
```

SetOpt sets the specified options passed as Env and returns an option to restore
the last set of arg's previous values

#### func (*Config) String

```go
func (c *Config) String() (s string)
```

String returns pretty print for the Blacklist struct

#### type Contenter

```go
type Contenter interface {
 Find(string) int
 GetList() *Objects
 Len() int
 SetURL(string, string)
 String() string
}
```

Contenter is an interface for handling the different file/http data sources

#### type Env

```go
type Env struct {

 // ioWriter io.Writer
 Log      *logging.Logger
 API      string        `json:"API,omitempty"`
 Arch     string        `json:"Arch,omitempty"`
 Bash     string        `json:"Bash,omitempty"`
 Cores    int           `json:"Cores,omitempty"`
 Disabled bool          `json:"Disabled"`
 Dbug     bool          `json:"Dbug,omitempty"`
 Dex      *list         `json:"Dex,omitempty"`
 Dir      string        `json:"Dir,omitempty"`
 DNSsvc   string        `json:"dnsmasq service,omitempty"`
 Exc      *list         `json:"Exc,omitempty"`
 Ext      string        `json:"dnsmasq fileExt.,omitempty"`
 File     string        `json:"File,omitempty"`
 FnFmt    string        `json:"File name fmt,omitempty"`
 InCLI    string        `json:"-"`
 Method   string        `json:"HTTP method,omitempty"`
 Pfx      dnsPfx        `json:"Prefix,omitempty"`
 Test     bool          `json:"Test,omitempty"`
 Timeout  time.Duration `json:"Timeout,omitempty"`
 Verb     bool          `json:"Verbosity,omitempty"`
 Wildcard `json:"Wildcard,omitempty"`
}
```

Env is struct of parameters

#### func (*Env) Debug

```go
func (e *Env) Debug(s ...interface{})
```

Debug logs debug messages when the Dbug flag is true

#### func (*Env) String

```go
func (e *Env) String() string
```

Env Stringer interface

#### type ExcDomnObjects

```go
type ExcDomnObjects struct {
 *Objects
}
```

ExcDomnObjects struct of *Objects for domain exclusions

#### func (*ExcDomnObjects) Find

```go
func (e *ExcDomnObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*ExcDomnObjects) GetList

```go
func (e *ExcDomnObjects) GetList() *Objects
```

GetList implements the Contenter interface for ExcDomnObjects

#### func (*ExcDomnObjects) Len

```go
func (e *ExcDomnObjects) Len() int
```

Len returns how many sources there are

#### func (*ExcDomnObjects) SetURL

```go
func (e *ExcDomnObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*ExcDomnObjects) String

```go
func (e *ExcDomnObjects) String() string
```

#### type ExcHostObjects

```go
type ExcHostObjects struct {
 *Objects
}
```

ExcHostObjects struct of *Objects for host exclusions

#### func (*ExcHostObjects) Find

```go
func (e *ExcHostObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*ExcHostObjects) GetList

```go
func (e *ExcHostObjects) GetList() *Objects
```

GetList implements the Contenter interface for ExcHostObjects

#### func (*ExcHostObjects) Len

```go
func (e *ExcHostObjects) Len() int
```

Len returns how many sources there are

#### func (*ExcHostObjects) SetURL

```go
func (e *ExcHostObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*ExcHostObjects) String

```go
func (e *ExcHostObjects) String() string
```

#### type ExcRootObjects

```go
type ExcRootObjects struct {
 *Objects
}
```

ExcRootObjects struct of *Objects for global domain exclusions

#### func (*ExcRootObjects) Find

```go
func (e *ExcRootObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*ExcRootObjects) GetList

```go
func (e *ExcRootObjects) GetList() *Objects
```

GetList implements the Contenter interface for ExcRootObjects

#### func (*ExcRootObjects) Len

```go
func (e *ExcRootObjects) Len() int
```

Len returns how many sources there are

#### func (*ExcRootObjects) SetURL

```go
func (e *ExcRootObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*ExcRootObjects) String

```go
func (e *ExcRootObjects) String() string
```

#### type FIODataObjects

```go
type FIODataObjects struct {
 *Objects
}
```

FIODataObjects struct of *Objects for files

#### func (*FIODataObjects) Find

```go
func (f *FIODataObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*FIODataObjects) GetList

```go
func (f *FIODataObjects) GetList() *Objects
```

GetList implements the Contenter interface for FIODataObjects

#### func (*FIODataObjects) Len

```go
func (f *FIODataObjects) Len() int
```

Len returns how many sources there are

#### func (*FIODataObjects) SetURL

```go
func (f *FIODataObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*FIODataObjects) String

```go
func (f *FIODataObjects) String() string
```

#### type FIODomnObjects

```go
type FIODomnObjects struct {
 *Objects
}
```

FIODomnObjects struct of *Objects for files

#### type FIOHostObjects

```go
type FIOHostObjects struct {
 *Objects
}
```

FIOHostObjects struct of *Objects for files

#### type IFace

```go
type IFace int
```

IFace type for labeling interface types

```go
const (
 Invalid IFace = iota + 100
 ExDmObj
 ExHtObj
 ExRtObj
 FileObj
 FylDObj
 FylHObj
 PreDObj
 PreHObj
 PreRObj
 URLdObj
 URLhObj
)
```

IFace types for labeling Content interfaces

#### func (IFace) String

```go
func (i IFace) String() string
```

#### type Objects

```go
type Objects struct {
 *Env
}
```

Objects is a struct of []*source

#### func (*Objects) Files

```go
func (o *Objects) Files() *CFile
```

Files returns a list of dnsmasq conf files from all srcs

#### func (*Objects) Filter

```go
func (o *Objects) Filter(ltype string) *Objects
```

Filter returns a subset of Objects filtered by ltype

#### func (*Objects) Find

```go
func (o *Objects) Find(elem string) int
```

Find returns the int position of an Objects' element

#### func (*Objects) Len

```go
func (o *Objects) Len() int
```

Implement Sort Interface for Objects

#### func (*Objects) Less

```go
func (o *Objects) Less(i, j int) bool
```

#### func (*Objects) Names

```go
func (o *Objects) Names() (s sort.StringSlice)
```

Names returns a sorted slice of Objects names

#### func (*Objects) String

```go
func (o *Objects) String() (s string)
```

Stringer for Objects

#### func (*Objects) Swap

```go
func (o *Objects) Swap(i, j int)
```

#### type Option

```go
type Option func(c *Config) Option
```

Option is a recursive function

#### func  API

```go
func API(s string) Option
```

API sets the EdgeOS CLI API command

#### func  Arch

```go
func Arch(s string) Option
```

Arch sets target CPU architecture

#### func  Bash

```go
func Bash(s string) Option
```

Bash sets the shell processor

#### func  Cores

```go
func Cores(i int) Option
```

Cores sets max CPU cores

#### func  DNSsvc

```go
func DNSsvc(s string) Option
```

DNSsvc sets dnsmasq restart command

#### func  Dbug

```go
func Dbug(b bool) Option
```

Dbug toggles Debug level on or off

#### func  Dir

```go
func Dir(s string) Option
```

Dir sets directory location

#### func  Disabled

```go
func Disabled(b bool) Option
```

Disabled toggles Disabled

#### func  Ext

```go
func Ext(s string) Option
```

Ext sets the blacklist file n extension

#### func  File

```go
func File(s string) Option
```

File sets the EdgeOS configuration file

#### func  FileNameFmt

```go
func FileNameFmt(s string) Option
```

FileNameFmt sets the EdgeOS configuration file name format

#### func  InCLI

```go
func InCLI(s string) Option
```

InCLI sets the CLI inSession command

#### func  Logger

```go
func Logger(l *logging.Logger) Option
```

Logger sets a pointer to the logger

#### func  Method

```go
func Method(s string) Option
```

Method sets the HTTP method

#### func  Prefix

```go
func Prefix(d string, h string) Option
```

Prefix sets the dnsmasq configuration address line prefix

#### func  Test

```go
func Test(b bool) Option
```

Test toggles testing mode on or off

#### func  Timeout

```go
func Timeout(t time.Duration) Option
```

Timeout sets how long before an unresponsive goroutine is aborted

#### func  Verb

```go
func Verb(b bool) Option
```

Verb sets the verbosity level to v

#### func  WCard

```go
func WCard(w Wildcard) Option
```

WCard sets file globbing wildcard values

#### type PreDomnObjects

```go
type PreDomnObjects struct {
 *Objects
}
```

PreDomnObjects struct of *Objects for pre-configured domains content

#### func (*PreDomnObjects) Find

```go
func (p *PreDomnObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*PreDomnObjects) GetList

```go
func (p *PreDomnObjects) GetList() *Objects
```

GetList implements the Contenter interface for PreDomnObjects

#### func (*PreDomnObjects) Len

```go
func (p *PreDomnObjects) Len() int
```

Len returns how many sources there are

#### func (*PreDomnObjects) SetURL

```go
func (p *PreDomnObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*PreDomnObjects) String

```go
func (p *PreDomnObjects) String() string
```

#### type PreHostObjects

```go
type PreHostObjects struct {
 *Objects
}
```

PreHostObjects struct of *Objects for pre-configured hosts content

#### func (*PreHostObjects) Find

```go
func (p *PreHostObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*PreHostObjects) GetList

```go
func (p *PreHostObjects) GetList() *Objects
```

GetList implements the Contenter interface for PreHostObjects

#### func (*PreHostObjects) Len

```go
func (p *PreHostObjects) Len() int
```

Len returns how many sources there are

#### func (*PreHostObjects) SetURL

```go
func (p *PreHostObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*PreHostObjects) String

```go
func (p *PreHostObjects) String() string
```

#### type PreRootObjects

```go
type PreRootObjects struct {
 *Objects
}
```

PreRootObjects struct of *Objects for pre-configured hosts content

#### func (*PreRootObjects) Find

```go
func (p *PreRootObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*PreRootObjects) GetList

```go
func (p *PreRootObjects) GetList() *Objects
```

GetList implements the Contenter interface for PreRootObjects

#### func (*PreRootObjects) Len

```go
func (p *PreRootObjects) Len() int
```

Len returns how many sources there are

#### func (*PreRootObjects) SetURL

```go
func (p *PreRootObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*PreRootObjects) String

```go
func (p *PreRootObjects) String() string
```

#### type URLDomnObjects

```go
type URLDomnObjects struct {
 *Objects
}
```

URLDomnObjects struct of *Objects for domain URLs

#### func (*URLDomnObjects) Find

```go
func (u *URLDomnObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*URLDomnObjects) GetList

```go
func (u *URLDomnObjects) GetList() *Objects
```

GetList implements the Contenter interface for URLDomnObjects

#### func (*URLDomnObjects) Len

```go
func (u *URLDomnObjects) Len() int
```

Len returns how many sources there are

#### func (*URLDomnObjects) SetURL

```go
func (u *URLDomnObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*URLDomnObjects) String

```go
func (u *URLDomnObjects) String() string
```

#### type URLHostObjects

```go
type URLHostObjects struct {
 *Objects
}
```

URLHostObjects struct of *Objects for host URLs

#### func (*URLHostObjects) Find

```go
func (u *URLHostObjects) Find(s string) int
```

Find returns the int position of an Objects' element

#### func (*URLHostObjects) GetList

```go
func (u *URLHostObjects) GetList() *Objects
```

GetList implements the Contenter interface for URLHostObjects

#### func (*URLHostObjects) Len

```go
func (u *URLHostObjects) Len() int
```

Len returns how many sources there are

#### func (*URLHostObjects) SetURL

```go
func (u *URLHostObjects) SetURL(name, url string)
```

SetURL sets the Object's url field value

#### func (*URLHostObjects) String

```go
func (u *URLHostObjects) String() string
```

#### type Wildcard

```go
type Wildcard struct {
 Node string `json:"Node,omitempty"`
 Name string `json:"Name,omitempty"`
}
```

Wildcard struct sets globbing wildcards for filename searches
