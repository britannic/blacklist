

# edgeos
`import "github.com/britannic/blacklist/internal/edgeos"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func GetFile(f string) (io.Reader, error)](#GetFile)
* [func Iter(i int) []struct{}](#Iter)
* [func NewWriter() io.Writer](#NewWriter)
* [type CFGcli](#CFGcli)
* [type CFGstatic](#CFGstatic)
* [type CFile](#CFile)
  * [func (c *CFile) Remove() error](#CFile.Remove)
  * [func (c *CFile) String() string](#CFile.String)
  * [func (c *CFile) Strings() []string](#CFile.Strings)
* [type ConfLoader](#ConfLoader)
* [type Config](#Config)
  * [func NewConfig(opts ...Option) *Config](#NewConfig)
  * [func (c *Config) Get(node string) *Objects](#Config.Get)
  * [func (c *Config) GetAll(ltypes ...string) *Objects](#Config.GetAll)
  * [func (c *Config) InSession() bool](#Config.InSession)
  * [func (c *Config) LTypes() []string](#Config.LTypes)
  * [func (c *Config) NewContent(iface IFace) (Contenter, error)](#Config.NewContent)
  * [func (c *Config) Nodes() (nodes []string)](#Config.Nodes)
  * [func (c *Config) ProcessContent(cts ...Contenter) error](#Config.ProcessContent)
  * [func (c *Config) ReadCfg(r ConfLoader) error](#Config.ReadCfg)
  * [func (c *Config) ReloadDNS() ([]byte, error)](#Config.ReloadDNS)
  * [func (c *Config) SetOpt(opts ...Option) Option](#Config.SetOpt)
  * [func (c *Config) String() (s string)](#Config.String)
* [type Contenter](#Contenter)
* [type ExcDomnObjects](#ExcDomnObjects)
  * [func (e *ExcDomnObjects) Find(elem string) int](#ExcDomnObjects.Find)
  * [func (e *ExcDomnObjects) GetList() *Objects](#ExcDomnObjects.GetList)
  * [func (e *ExcDomnObjects) Len() int](#ExcDomnObjects.Len)
  * [func (e *ExcDomnObjects) SetURL(name, url string)](#ExcDomnObjects.SetURL)
  * [func (e *ExcDomnObjects) String() string](#ExcDomnObjects.String)
* [type ExcHostObjects](#ExcHostObjects)
  * [func (e *ExcHostObjects) Find(elem string) int](#ExcHostObjects.Find)
  * [func (e *ExcHostObjects) GetList() *Objects](#ExcHostObjects.GetList)
  * [func (e *ExcHostObjects) Len() int](#ExcHostObjects.Len)
  * [func (e *ExcHostObjects) SetURL(name, url string)](#ExcHostObjects.SetURL)
  * [func (e *ExcHostObjects) String() string](#ExcHostObjects.String)
* [type ExcRootObjects](#ExcRootObjects)
  * [func (e *ExcRootObjects) Find(elem string) int](#ExcRootObjects.Find)
  * [func (e *ExcRootObjects) GetList() *Objects](#ExcRootObjects.GetList)
  * [func (e *ExcRootObjects) Len() int](#ExcRootObjects.Len)
  * [func (e *ExcRootObjects) SetURL(name, url string)](#ExcRootObjects.SetURL)
  * [func (e *ExcRootObjects) String() string](#ExcRootObjects.String)
* [type FIODataObjects](#FIODataObjects)
  * [func (f *FIODataObjects) Find(elem string) int](#FIODataObjects.Find)
  * [func (f *FIODataObjects) GetList() *Objects](#FIODataObjects.GetList)
  * [func (f *FIODataObjects) Len() int](#FIODataObjects.Len)
  * [func (f *FIODataObjects) SetURL(name, url string)](#FIODataObjects.SetURL)
  * [func (f *FIODataObjects) String() string](#FIODataObjects.String)
* [type IFace](#IFace)
  * [func (i IFace) String() (s string)](#IFace.String)
* [type Msg](#Msg)
  * [func NewMsg(s string) *Msg](#NewMsg)
  * [func (m *Msg) GetTotal() int](#Msg.GetTotal)
  * [func (m *Msg) ReadDupes() int](#Msg.ReadDupes)
  * [func (m *Msg) ReadNew() int](#Msg.ReadNew)
  * [func (m *Msg) ReadUniq() int](#Msg.ReadUniq)
  * [func (m *Msg) String() string](#Msg.String)
* [type Objects](#Objects)
  * [func (o *Objects) Files() *CFile](#Objects.Files)
  * [func (o *Objects) Filter(ltype string) *Objects](#Objects.Filter)
  * [func (o *Objects) Find(elem string) int](#Objects.Find)
  * [func (o *Objects) Len() int](#Objects.Len)
  * [func (o *Objects) Less(i, j int) bool](#Objects.Less)
  * [func (o *Objects) Names() (s sort.StringSlice)](#Objects.Names)
  * [func (o *Objects) String() string](#Objects.String)
  * [func (o *Objects) Swap(i, j int)](#Objects.Swap)
* [type Option](#Option)
  * [func API(s string) Option](#API)
  * [func Arch(arch string) Option](#Arch)
  * [func Bash(cmd string) Option](#Bash)
  * [func Cores(i int) Option](#Cores)
  * [func DNSsvc(d string) Option](#DNSsvc)
  * [func Dbug(b bool) Option](#Dbug)
  * [func Dir(d string) Option](#Dir)
  * [func Disabled(b bool) Option](#Disabled)
  * [func Ext(e string) Option](#Ext)
  * [func File(f string) Option](#File)
  * [func FileNameFmt(f string) Option](#FileNameFmt)
  * [func InCLI(in string) Option](#InCLI)
  * [func LTypes(s []string) Option](#LTypes)
  * [func Level(s string) Option](#Level)
  * [func Logger(l *logging.Logger) Option](#Logger)
  * [func Method(method string) Option](#Method)
  * [func Prefix(l string) Option](#Prefix)
  * [func Test(b bool) Option](#Test)
  * [func Timeout(t time.Duration) Option](#Timeout)
  * [func Verb(b bool) Option](#Verb)
  * [func WCard(w Wildcard) Option](#WCard)
  * [func Writer(w io.Writer) Option](#Writer)
* [type Parms](#Parms)
  * [func (p *Parms) String() string](#Parms.String)
* [type PreDomnObjects](#PreDomnObjects)
  * [func (p *PreDomnObjects) Find(elem string) int](#PreDomnObjects.Find)
  * [func (p *PreDomnObjects) GetList() *Objects](#PreDomnObjects.GetList)
  * [func (p *PreDomnObjects) Len() int](#PreDomnObjects.Len)
  * [func (p *PreDomnObjects) SetURL(name, url string)](#PreDomnObjects.SetURL)
  * [func (p *PreDomnObjects) String() string](#PreDomnObjects.String)
* [type PreHostObjects](#PreHostObjects)
  * [func (p *PreHostObjects) Find(elem string) int](#PreHostObjects.Find)
  * [func (p *PreHostObjects) GetList() *Objects](#PreHostObjects.GetList)
  * [func (p *PreHostObjects) Len() int](#PreHostObjects.Len)
  * [func (p *PreHostObjects) SetURL(name, url string)](#PreHostObjects.SetURL)
  * [func (p *PreHostObjects) String() string](#PreHostObjects.String)
* [type URLDomnObjects](#URLDomnObjects)
  * [func (u *URLDomnObjects) Find(elem string) int](#URLDomnObjects.Find)
  * [func (u *URLDomnObjects) GetList() *Objects](#URLDomnObjects.GetList)
  * [func (u *URLDomnObjects) Len() int](#URLDomnObjects.Len)
  * [func (u *URLDomnObjects) SetURL(name, url string)](#URLDomnObjects.SetURL)
  * [func (u *URLDomnObjects) String() string](#URLDomnObjects.String)
* [type URLHostObjects](#URLHostObjects)
  * [func (u *URLHostObjects) Find(elem string) int](#URLHostObjects.Find)
  * [func (u *URLHostObjects) GetList() *Objects](#URLHostObjects.GetList)
  * [func (u *URLHostObjects) Len() int](#URLHostObjects.Len)
  * [func (u *URLHostObjects) SetURL(name, url string)](#URLHostObjects.SetURL)
  * [func (u *URLHostObjects) String() string](#URLHostObjects.String)
* [type Wildcard](#Wildcard)


#### <a name="pkg-files">Package files</a>
[config.go](/src/github.com/britannic/blacklist/internal/edgeos/config.go) [content.go](/src/github.com/britannic/blacklist/internal/edgeos/content.go) [data.go](/src/github.com/britannic/blacklist/internal/edgeos/data.go) [http.go](/src/github.com/britannic/blacklist/internal/edgeos/http.go) [io.go](/src/github.com/britannic/blacklist/internal/edgeos/io.go) [json.go](/src/github.com/britannic/blacklist/internal/edgeos/json.go) [list.go](/src/github.com/britannic/blacklist/internal/edgeos/list.go) [msg.go](/src/github.com/britannic/blacklist/internal/edgeos/msg.go) [ntype_string.go](/src/github.com/britannic/blacklist/internal/edgeos/ntype_string.go) [object.go](/src/github.com/britannic/blacklist/internal/edgeos/object.go) [opts.go](/src/github.com/britannic/blacklist/internal/edgeos/opts.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (

    // ExcDomns labels domain exclusions
    ExcDomns = "excluded-domains"
    // ExcHosts labels host exclusions
    ExcHosts = "excluded-hosts"
    // ExcRoots labels global domain exclusions
    ExcRoots = "excluded-global"
    // PreDomns designates string label for preconfigured blacklisted domains
    PreDomns = "domains." + preNoun
    // PreHosts designates string label for preconfigured blacklisted hosts
    PreHosts = "hosts." + preNoun
    // False is a string constant
    False = "false"
    // True is a string constant
    True = "true"
)
```



## <a name="GetFile">func</a> [GetFile](/src/target/io.go?s=1940:1981#L83)
``` go
func GetFile(f string) (io.Reader, error)
```
GetFile reads a file and returns an io.Reader



## <a name="Iter">func</a> [Iter](/src/target/data.go?s=2429:2456#L100)
``` go
func Iter(i int) []struct{}
```
Iter iterates over ints - use it in for loops



## <a name="NewWriter">func</a> [NewWriter](/src/target/data.go?s=2524:2550#L105)
``` go
func NewWriter() io.Writer
```
NewWriter returns an io.Writer




## <a name="CFGcli">type</a> [CFGcli](/src/target/io.go?s=139:182#L4)
``` go
type CFGcli struct {
    *Config
    Cfg string
}
```
CFGcli loads configurations using the EdgeOS CFGcli










## <a name="CFGstatic">type</a> [CFGstatic](/src/target/io.go?s=237:283#L10)
``` go
type CFGstatic struct {
    *Config
    Cfg string
}
```
CFGstatic loads static configurations for testing










## <a name="CFile">type</a> [CFile](/src/target/config.go?s=493:551#L18)
``` go
type CFile struct {
    *Parms
    Names []string
    // contains filtered or unexported fields
}
```
CFile holds an array of file names










### <a name="CFile.Remove">func</a> (\*CFile) [Remove](/src/target/config.go?s=8695:8725#L401)
``` go
func (c *CFile) Remove() error
```
Remove deletes a CFile array of file names




### <a name="CFile.String">func</a> (\*CFile) [String](/src/target/config.go?s=10103:10134#L457)
``` go
func (c *CFile) String() string
```
String implements string method




### <a name="CFile.Strings">func</a> (\*CFile) [Strings](/src/target/config.go?s=10222:10256#L462)
``` go
func (c *CFile) Strings() []string
```
Strings returns a sorted array of strings.




## <a name="ConfLoader">type</a> [ConfLoader](/src/target/config.go?s=406:453#L13)
``` go
type ConfLoader interface {
    // contains filtered or unexported methods
}
```
ConfLoader interface defines configuration load method










## <a name="Config">type</a> [Config](/src/target/config.go?s=599:635#L25)
``` go
type Config struct {
    *Parms
    // contains filtered or unexported fields
}
```
Config is a struct of configuration fields







### <a name="NewConfig">func</a> [NewConfig](/src/target/opts.go?s=5178:5216#L229)
``` go
func NewConfig(opts ...Option) *Config
```
NewConfig returns a new *Config initialized with the parameter options passed to it





### <a name="Config.Get">func</a> (\*Config) [Get](/src/target/config.go?s=3950:3992#L191)
``` go
func (c *Config) Get(node string) *Objects
```
Get returns an *Object for a given node




### <a name="Config.GetAll">func</a> (\*Config) [GetAll](/src/target/config.go?s=4271:4321#L210)
``` go
func (c *Config) GetAll(ltypes ...string) *Objects
```
GetAll returns an array of Objects




### <a name="Config.InSession">func</a> (\*Config) [InSession](/src/target/config.go?s=5115:5148#L253)
``` go
func (c *Config) InSession() bool
```
InSession returns true if VyOS/EdgeOS configuration is in session




### <a name="Config.LTypes">func</a> (\*Config) [LTypes](/src/target/config.go?s=10348:10382#L468)
``` go
func (c *Config) LTypes() []string
```
LTypes returns an array of configured nodes




### <a name="Config.NewContent">func</a> (\*Config) [NewContent](/src/target/config.go?s=2563:2622#L125)
``` go
func (c *Config) NewContent(iface IFace) (Contenter, error)
```
NewContent returns an interface of the requested IFace type




### <a name="Config.Nodes">func</a> (\*Config) [Nodes](/src/target/config.go?s=5536:5577#L266)
``` go
func (c *Config) Nodes() (nodes []string)
```
Nodes returns an array of configured nodes




### <a name="Config.ProcessContent">func</a> (\*Config) [ProcessContent](/src/target/content.go?s=7939:7994#L380)
``` go
func (c *Config) ProcessContent(cts ...Contenter) error
```
ProcessContent processes the Contents array




### <a name="Config.ReadCfg">func</a> (\*Config) [ReadCfg](/src/target/config.go?s=5931:5975#L285)
``` go
func (c *Config) ReadCfg(r ConfLoader) error
```
ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure




### <a name="Config.ReloadDNS">func</a> (\*Config) [ReloadDNS](/src/target/config.go?s=8500:8544#L394)
``` go
func (c *Config) ReloadDNS() ([]byte, error)
```
ReloadDNS reloads the dnsmasq configuration




### <a name="Config.SetOpt">func</a> (\*Config) [SetOpt](/src/target/opts.go?s=2018:2064#L68)
``` go
func (c *Config) SetOpt(opts ...Option) Option
```
SetOpt sets the specified options passed as Parms and returns an option to restore the last set of arg's previous values




### <a name="Config.String">func</a> (\*Config) [String](/src/target/config.go?s=9210:9246#L422)
``` go
func (c *Config) String() (s string)
```
String returns pretty print for the Blacklist struct




## <a name="Contenter">type</a> [Contenter](/src/target/content.go?s=445:563#L28)
``` go
type Contenter interface {
    Find(string) int
    GetList() *Objects
    Len() int
    SetURL(string, string)
    String() string
}
```
Contenter is a Content interface










## <a name="ExcDomnObjects">type</a> [ExcDomnObjects](/src/target/content.go?s=624:664#L37)
``` go
type ExcDomnObjects struct {
    *Objects
}
```
ExcDomnObjects implements GetList for domain exclusions










### <a name="ExcDomnObjects.Find">func</a> (\*ExcDomnObjects) [Find](/src/target/content.go?s=1420:1466#L77)
``` go
func (e *ExcDomnObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="ExcDomnObjects.GetList">func</a> (\*ExcDomnObjects) [GetList](/src/target/content.go?s=2909:2952#L157)
``` go
func (e *ExcDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcDomnObjects




### <a name="ExcDomnObjects.Len">func</a> (\*ExcDomnObjects) [Len](/src/target/content.go?s=5474:5508#L285)
``` go
func (e *ExcDomnObjects) Len() int
```
Len returns how many objects there are




### <a name="ExcDomnObjects.SetURL">func</a> (\*ExcDomnObjects) [SetURL](/src/target/content.go?s=8768:8817#L424)
``` go
func (e *ExcDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcDomnObjects.String">func</a> (\*ExcDomnObjects) [String](/src/target/content.go?s=10068:10108#L495)
``` go
func (e *ExcDomnObjects) String() string
```



## <a name="ExcHostObjects">type</a> [ExcHostObjects](/src/target/content.go?s=723:763#L42)
``` go
type ExcHostObjects struct {
    *Objects
}
```
ExcHostObjects implements GetList for host exclusions










### <a name="ExcHostObjects.Find">func</a> (\*ExcHostObjects) [Find](/src/target/content.go?s=1605:1651#L87)
``` go
func (e *ExcHostObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="ExcHostObjects.GetList">func</a> (\*ExcHostObjects) [GetList](/src/target/content.go?s=3187:3230#L171)
``` go
func (e *ExcHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcHostObjects




### <a name="ExcHostObjects.Len">func</a> (\*ExcHostObjects) [Len](/src/target/content.go?s=5580:5614#L288)
``` go
func (e *ExcHostObjects) Len() int
```
Len returns how many objects there are




### <a name="ExcHostObjects.SetURL">func</a> (\*ExcHostObjects) [SetURL](/src/target/content.go?s=8936:8985#L433)
``` go
func (e *ExcHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcHostObjects.String">func</a> (\*ExcHostObjects) [String](/src/target/content.go?s=10139:10179#L496)
``` go
func (e *ExcHostObjects) String() string
```



## <a name="ExcRootObjects">type</a> [ExcRootObjects](/src/target/content.go?s=831:871#L47)
``` go
type ExcRootObjects struct {
    *Objects
}
```
ExcRootObjects implements GetList for global domain exclusions










### <a name="ExcRootObjects.Find">func</a> (\*ExcRootObjects) [Find](/src/target/content.go?s=1790:1836#L97)
``` go
func (e *ExcRootObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="ExcRootObjects.GetList">func</a> (\*ExcRootObjects) [GetList](/src/target/content.go?s=3465:3508#L185)
``` go
func (e *ExcRootObjects) GetList() *Objects
```
GetList implements the Contenter interface for ExcRootObjects




### <a name="ExcRootObjects.Len">func</a> (\*ExcRootObjects) [Len](/src/target/content.go?s=5686:5720#L291)
``` go
func (e *ExcRootObjects) Len() int
```
Len returns how many objects there are




### <a name="ExcRootObjects.SetURL">func</a> (\*ExcRootObjects) [SetURL](/src/target/content.go?s=9104:9153#L442)
``` go
func (e *ExcRootObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="ExcRootObjects.String">func</a> (\*ExcRootObjects) [String](/src/target/content.go?s=10210:10250#L497)
``` go
func (e *ExcRootObjects) String() string
```



## <a name="FIODataObjects">type</a> [FIODataObjects](/src/target/content.go?s=920:960#L52)
``` go
type FIODataObjects struct {
    *Objects
}
```
FIODataObjects implements GetList for files










### <a name="FIODataObjects.Find">func</a> (\*FIODataObjects) [Find](/src/target/content.go?s=1975:2021#L107)
``` go
func (f *FIODataObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="FIODataObjects.GetList">func</a> (\*FIODataObjects) [GetList](/src/target/content.go?s=3743:3786#L199)
``` go
func (f *FIODataObjects) GetList() *Objects
```
GetList implements the Contenter interface for FIODataObjects




### <a name="FIODataObjects.Len">func</a> (\*FIODataObjects) [Len](/src/target/content.go?s=5792:5826#L294)
``` go
func (f *FIODataObjects) Len() int
```
Len returns how many objects there are




### <a name="FIODataObjects.SetURL">func</a> (\*FIODataObjects) [SetURL](/src/target/content.go?s=9272:9321#L451)
``` go
func (f *FIODataObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="FIODataObjects.String">func</a> (\*FIODataObjects) [String](/src/target/content.go?s=10281:10321#L498)
``` go
func (f *FIODataObjects) String() string
```



## <a name="IFace">type</a> [IFace](/src/target/content.go?s=186:200#L7)
``` go
type IFace int
```
IFace type for labeling interface types


``` go
const (
    Invalid IFace = iota + 100
    ExRtObj
    ExDmObj
    ExHtObj
    FileObj
    PreDObj
    PreHObj
    URLdObj
    URLhObj
)
```
IFace types for labeling interface types










### <a name="IFace.String">func</a> (IFace) [String](/src/target/content.go?s=10637:10671#L504)
``` go
func (i IFace) String() (s string)
```



## <a name="Msg">type</a> [Msg](/src/target/msg.go?s=112:301#L1)
``` go
type Msg struct {
    Name string `json:"name"`
    Done bool   `json:"done"`
    *sync.RWMutex
    Dupes int `json:"dupes"`
    New   int `json:"new"`
    Total int `json:"total"`
    Uniq  int `json:"uniq"`
}
```
Msg is a struct for recording stats from ProcessContent







### <a name="NewMsg">func</a> [NewMsg](/src/target/msg.go?s=760:786#L39)
``` go
func NewMsg(s string) *Msg
```
NewMsg initializes and returns a new Msg struct





### <a name="Msg.GetTotal">func</a> (\*Msg) [GetTotal](/src/target/msg.go?s=343:371#L10)
``` go
func (m *Msg) GetTotal() int
```
GetTotal returns total m.New records




### <a name="Msg.ReadDupes">func</a> (\*Msg) [ReadDupes](/src/target/msg.go?s=899:928#L47)
``` go
func (m *Msg) ReadDupes() int
```
ReadDupes returns current value of Msg.Dupes




### <a name="Msg.ReadNew">func</a> (\*Msg) [ReadNew](/src/target/msg.go?s=1024:1051#L54)
``` go
func (m *Msg) ReadNew() int
```
ReadNew returns current value of Msg.New




### <a name="Msg.ReadUniq">func</a> (\*Msg) [ReadUniq](/src/target/msg.go?s=1147:1175#L61)
``` go
func (m *Msg) ReadUniq() int
```
ReadUniq returns current value of msg.Uniq




### <a name="Msg.String">func</a> (\*Msg) [String](/src/target/msg.go?s=1226:1255#L67)
``` go
func (m *Msg) String() string
```



## <a name="Objects">type</a> [Objects](/src/target/object.go?s=401:445#L20)
``` go
type Objects struct {
    *Parms
    // contains filtered or unexported fields
}
```
Objects is a struct of []*Object










### <a name="Objects.Files">func</a> (\*Objects) [Files](/src/target/object.go?s=1115:1147#L52)
``` go
func (o *Objects) Files() *CFile
```
Files returns a list of dnsmasq conf files from all srcs




### <a name="Objects.Filter">func</a> (\*Objects) [Filter](/src/target/object.go?s=1501:1548#L66)
``` go
func (o *Objects) Filter(ltype string) *Objects
```
Filter returns a subset of Objects; ltypes with "-" prepended remove ltype




### <a name="Objects.Find">func</a> (\*Objects) [Find](/src/target/object.go?s=2283:2322#L105)
``` go
func (o *Objects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="Objects.Len">func</a> (\*Objects) [Len](/src/target/object.go?s=3659:3686#L166)
``` go
func (o *Objects) Len() int
```
Implement Sort Interface for Objects




### <a name="Objects.Less">func</a> (\*Objects) [Less](/src/target/object.go?s=3717:3754#L167)
``` go
func (o *Objects) Less(i, j int) bool
```



### <a name="Objects.Names">func</a> (\*Objects) [Names](/src/target/object.go?s=2753:2799#L129)
``` go
func (o *Objects) Names() (s sort.StringSlice)
```
Names returns a sorted slice of Objects names




### <a name="Objects.String">func</a> (\*Objects) [String](/src/target/object.go?s=3556:3589#L161)
``` go
func (o *Objects) String() string
```
Stringer for Objects




### <a name="Objects.Swap">func</a> (\*Objects) [Swap](/src/target/object.go?s=3792:3824#L168)
``` go
func (o *Objects) Swap(i, j int)
```



## <a name="Option">type</a> [Option](/src/target/opts.go?s=1628:1662#L47)
``` go
type Option func(c *Config) Option
```
Option is a recursive function







### <a name="API">func</a> [API](/src/target/opts.go?s=2550:2575#L93)
``` go
func API(s string) Option
```
API sets the EdgeOS CLI API command


### <a name="Arch">func</a> [Arch](/src/target/opts.go?s=2379:2408#L84)
``` go
func Arch(arch string) Option
```
Arch sets target CPU architecture


### <a name="Bash">func</a> [Bash](/src/target/opts.go?s=2705:2733#L102)
``` go
func Bash(cmd string) Option
```
Bash sets the shell processor


### <a name="Cores">func</a> [Cores](/src/target/opts.go?s=2863:2887#L111)
``` go
func Cores(i int) Option
```
Cores sets max CPU cores


### <a name="DNSsvc">func</a> [DNSsvc](/src/target/opts.go?s=3537:3565#L148)
``` go
func DNSsvc(d string) Option
```
DNSsvc sets dnsmasq restart command


### <a name="Dbug">func</a> [Dbug](/src/target/opts.go?s=3221:3245#L130)
``` go
func Dbug(b bool) Option
```
Dbug toggles debug level on or off


### <a name="Dir">func</a> [Dir](/src/target/opts.go?s=3376:3401#L139)
``` go
func Dir(d string) Option
```
Dir sets directory location


### <a name="Disabled">func</a> [Disabled](/src/target/opts.go?s=3043:3071#L121)
``` go
func Disabled(b bool) Option
```
Disabled toggles Disabled


### <a name="Ext">func</a> [Ext](/src/target/opts.go?s=3714:3739#L157)
``` go
func Ext(e string) Option
```
Ext sets the blacklist file n extension


### <a name="File">func</a> [File](/src/target/opts.go?s=3879:3905#L166)
``` go
func File(f string) Option
```
File sets the EdgeOS configuration file


### <a name="FileNameFmt">func</a> [FileNameFmt](/src/target/opts.go?s=4067:4100#L175)
``` go
func FileNameFmt(f string) Option
```
FileNameFmt sets the EdgeOS configuration file name format


### <a name="InCLI">func</a> [InCLI](/src/target/opts.go?s=4249:4277#L184)
``` go
func InCLI(in string) Option
```
InCLI sets the CLI inSession command


### <a name="LTypes">func</a> [LTypes](/src/target/opts.go?s=4780:4810#L211)
``` go
func LTypes(s []string) Option
```
LTypes sets an array of legal types used by Source


### <a name="Level">func</a> [Level](/src/target/opts.go?s=4420:4447#L193)
``` go
func Level(s string) Option
```
Level sets the EdgeOS API CLI level


### <a name="Logger">func</a> [Logger](/src/target/opts.go?s=4589:4626#L202)
``` go
func Logger(l *logging.Logger) Option
```
Logger sets a pointer to the logger


### <a name="Method">func</a> [Method](/src/target/opts.go?s=4947:4980#L220)
``` go
func Method(method string) Option
```
Method sets the HTTP method


### <a name="Prefix">func</a> [Prefix](/src/target/opts.go?s=5543:5571#L245)
``` go
func Prefix(l string) Option
```
Prefix sets the dnsmasq configuration address line prefix


### <a name="Test">func</a> [Test](/src/target/opts.go?s=5860:5884#L260)
``` go
func Test(b bool) Option
```
Test toggles testing mode on or off


### <a name="Timeout">func</a> [Timeout](/src/target/opts.go?s=6053:6089#L269)
``` go
func Timeout(t time.Duration) Option
```
Timeout sets how long before an unresponsive goroutine is aborted


### <a name="Verb">func</a> [Verb](/src/target/opts.go?s=6236:6260#L278)
``` go
func Verb(b bool) Option
```
Verb sets the verbosity level to v


### <a name="WCard">func</a> [WCard](/src/target/opts.go?s=6404:6433#L287)
``` go
func WCard(w Wildcard) Option
```
WCard sets file globbing wildcard values


### <a name="Writer">func</a> [Writer](/src/target/opts.go?s=6608:6639#L296)
``` go
func Writer(w io.Writer) Option
```
Writer provides an address for anything that can use io.Writer





## <a name="Parms">type</a> [Parms](/src/target/opts.go?s=148:1330#L4)
``` go
type Parms struct {
    Log      *logging.Logger
    API      string        `json:"API,omitempty"`
    Arch     string        `json:"Arch,omitempty"`
    Bash     string        `json:"Bash,omitempty"`
    Cores    int           `json:"Cores,omitempty"`
    Disabled bool          `json:"Disabled"`
    Dbug     bool          `json:"Dbug,omitempty"`
    Dex      list          `json:"Dex,omitempty"`
    Dir      string        `json:"Dir,omitempty"`
    DNSsvc   string        `json:"dnsmasq service,omitempty"`
    Exc      list          `json:"Exc,omitempty"`
    Ext      string        `json:"dnsmasq fileExt.,omitempty"`
    File     string        `json:"File,omitempty"`
    FnFmt    string        `json:"File name fmt,omitempty"`
    InCLI    string        `json:"-"`
    Level    string        `json:"CLI Path,omitempty"`
    Ltypes   []string      `json:"Leaf nodes,omitempty"`
    Method   string        `json:"HTTP method,omitempty"`
    Pfx      string        `json:"Prefix,omitempty"`
    Test     bool          `json:"Test,omitempty"`
    Timeout  time.Duration `json:"Timeout,omitempty"`
    Verb     bool          `json:"Verbosity,omitempty"`
    Wildcard `json:"Wildcard,omitempty"`
    // contains filtered or unexported fields
}
```
Parms is struct of parameters










### <a name="Parms.String">func</a> (\*Parms) [String](/src/target/opts.go?s=5721:5752#L254)
``` go
func (p *Parms) String() string
```
String method to implement fmt.Print interface




## <a name="PreDomnObjects">type</a> [PreDomnObjects](/src/target/content.go?s=1034:1074#L57)
``` go
type PreDomnObjects struct {
    *Objects
}
```
PreDomnObjects implements GetList for pre-configured domains content










### <a name="PreDomnObjects.Find">func</a> (\*PreDomnObjects) [Find](/src/target/content.go?s=2160:2206#L117)
``` go
func (p *PreDomnObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="PreDomnObjects.GetList">func</a> (\*PreDomnObjects) [GetList](/src/target/content.go?s=4170:4213#L220)
``` go
func (p *PreDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreDomnObjects




### <a name="PreDomnObjects.Len">func</a> (\*PreDomnObjects) [Len](/src/target/content.go?s=5898:5932#L297)
``` go
func (p *PreDomnObjects) Len() int
```
Len returns how many objects there are




### <a name="PreDomnObjects.SetURL">func</a> (\*PreDomnObjects) [SetURL](/src/target/content.go?s=9440:9489#L460)
``` go
func (p *PreDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="PreDomnObjects.String">func</a> (\*PreDomnObjects) [String](/src/target/content.go?s=10352:10392#L499)
``` go
func (p *PreDomnObjects) String() string
```



## <a name="PreHostObjects">type</a> [PreHostObjects](/src/target/content.go?s=1146:1186#L62)
``` go
type PreHostObjects struct {
    *Objects
}
```
PreHostObjects implements GetList for pre-configured hosts content










### <a name="PreHostObjects.Find">func</a> (\*PreHostObjects) [Find](/src/target/content.go?s=2345:2391#L127)
``` go
func (p *PreHostObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="PreHostObjects.GetList">func</a> (\*PreHostObjects) [GetList](/src/target/content.go?s=4428:4471#L231)
``` go
func (p *PreHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for PreHostObjects




### <a name="PreHostObjects.Len">func</a> (\*PreHostObjects) [Len](/src/target/content.go?s=6004:6038#L300)
``` go
func (p *PreHostObjects) Len() int
```
Len returns how many objects there are




### <a name="PreHostObjects.SetURL">func</a> (\*PreHostObjects) [SetURL](/src/target/content.go?s=9608:9657#L469)
``` go
func (p *PreHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="PreHostObjects.String">func</a> (\*PreHostObjects) [String](/src/target/content.go?s=10423:10463#L500)
``` go
func (p *PreHostObjects) String() string
```



## <a name="URLDomnObjects">type</a> [URLDomnObjects](/src/target/content.go?s=1234:1274#L67)
``` go
type URLDomnObjects struct {
    *Objects
}
```
URLDomnObjects implements GetList for URLs










### <a name="URLDomnObjects.Find">func</a> (\*URLDomnObjects) [Find](/src/target/content.go?s=2715:2761#L147)
``` go
func (u *URLDomnObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="URLDomnObjects.GetList">func</a> (\*URLDomnObjects) [GetList](/src/target/content.go?s=4686:4729#L242)
``` go
func (u *URLDomnObjects) GetList() *Objects
```
GetList implements the Contenter interface for URLHostObjects




### <a name="URLDomnObjects.Len">func</a> (\*URLDomnObjects) [Len](/src/target/content.go?s=6110:6144#L303)
``` go
func (u *URLDomnObjects) Len() int
```
Len returns how many objects there are




### <a name="URLDomnObjects.SetURL">func</a> (\*URLDomnObjects) [SetURL](/src/target/content.go?s=9776:9825#L478)
``` go
func (u *URLDomnObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="URLDomnObjects.String">func</a> (\*URLDomnObjects) [String](/src/target/content.go?s=10494:10534#L501)
``` go
func (u *URLDomnObjects) String() string
```



## <a name="URLHostObjects">type</a> [URLHostObjects](/src/target/content.go?s=1322:1362#L72)
``` go
type URLHostObjects struct {
    *Objects
}
```
URLHostObjects implements GetList for URLs










### <a name="URLHostObjects.Find">func</a> (\*URLHostObjects) [Find](/src/target/content.go?s=2530:2576#L137)
``` go
func (u *URLHostObjects) Find(elem string) int
```
Find returns the int position of an Objects' element




### <a name="URLHostObjects.GetList">func</a> (\*URLHostObjects) [GetList](/src/target/content.go?s=5091:5134#L263)
``` go
func (u *URLHostObjects) GetList() *Objects
```
GetList implements the Contenter interface for URLHostObjects




### <a name="URLHostObjects.Len">func</a> (\*URLHostObjects) [Len](/src/target/content.go?s=6216:6250#L306)
``` go
func (u *URLHostObjects) Len() int
```
Len returns how many objects there are




### <a name="URLHostObjects.SetURL">func</a> (\*URLHostObjects) [SetURL](/src/target/content.go?s=9944:9993#L487)
``` go
func (u *URLHostObjects) SetURL(name, url string)
```
SetURL sets the Object's url field value




### <a name="URLHostObjects.String">func</a> (\*URLHostObjects) [String](/src/target/content.go?s=10565:10605#L502)
``` go
func (u *URLHostObjects) String() string
```



## <a name="Wildcard">type</a> [Wildcard](/src/target/opts.go?s=1494:1592#L41)
``` go
type Wildcard struct {
    Node string `json:"Node,omitempty"`
    Name string `json:"Name,omitempty"`
}
```
Wildcard struct sets globbing wildcards for filename searches














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
