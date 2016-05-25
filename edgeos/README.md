
# edgeos
    import "github.com/britannic/blacklist/edgeos"

Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.




## Constants
``` go
const (

    // Fext sets the dnsmasq configuration file extension
    Fext = "blacklist.conf"

    // Domains sets the domains string
    Domains = "domains"

    // Hosts sets the hosts string
    Hosts = "hosts"

    // PreConf sets the string for pre-configured
    PreConf = "pre-configured"

    // Root is the topmost node
    Root = blacklist

    // False is a string constant
    False = "false"

    // True is a string constant
    True = "true"
)
```

## Variables
``` go
var (
    // API sets the path and executable for the EdgeOS shell API
    API = "/bin/cli-shell-api"
)
```

## func APICmd
``` go
func APICmd() (r map[string]string)
```
APICmd returns a map of CLI commands


## func DeleteFile
``` go
func DeleteFile(f string) bool
```
DeleteFile removes a file if it exists


## func DiffArray
``` go
func DiffArray(a, b []string) (diff []string)
```
DiffArray returns the delta of two arrays


## func GetHTTP
``` go
func GetHTTP(method, URL string) (io.Reader, error)
```
GetHTTP creates http requests to download data


## func getType
``` go
func getType(in interface{}) (out interface{})
```
getType returns the converted "in" type


## func Insession
``` go
func Insession() bool
```
Insession returns true if VyOS/EdgeOS configuration is in session


## func ListFiles
``` go
func ListFiles(dir string) (files []string, err error)
```
ListFiles returns a list of blacklist files


## func Load
``` go
func Load(action, level string) (reader io.Reader, err error)
```
Load reads the config using the EdgeOS/VyOS cli-shell-api


## func PurgeFiles
``` go
func PurgeFiles(files []string) (err error)
```
PurgeFiles removes any orphaned blacklist files that don't have sources


## func SHCmd
``` go
func SHCmd(a string) string
```
SHCmd returns the appropriate command for non-tty or tty configure context


## func ToBool
``` go
func ToBool(s string) bool
```
ToBool converts a string ("true" or "false") to it's boolean equivalent


## func WriteFile
``` go
func WriteFile(fname string, data io.Reader) (err error)
```
WriteFile writes blacklist data to storage



## type Config
``` go
type Config map[string]*EdgeOS
```
Config is a map of EdgeOS











### func (Config) Disabled
``` go
func (c *Config) Disabled(node string) bool
```
Disabled returns the node is true or false



### func (Config) Excludes
``` go
func (c *Config) Excludes(node string) []string
```
Excludes returns an array of excluded blacklist domains/hosts



### func (Config) Files
``` go
func (c *Config) Files(dir string, nodes []string) (files []string)
```
Files returns a list of dnsmasq conf files from all srcs



### func (Config) FormatData
``` go
func (c *Config) FormatData(fmttr string, data []string) (reader io.Reader, list List)
```
FormatData returns a io.Reader loaded with dnsmasq formatted data



### func (Config) Get
``` go
func (c *Config) Get(node string) (e *EdgeOS)
```
Get returns a normalized EdgeOS data set



### func (Config) GetExcludes
``` go
func (c *Config) GetExcludes(dex, ex List, nodes []string) (List, List)
```
GetExcludes collates the configured excludes and merges the ex/dex lists



### func (Config) IP
``` go
func (c *Config) IP(node string) string
```
IP returns the configured node IP, or the root node's IP if ""



### func (Config) Includes
``` go
func (c *Config) Includes(node string) []string
```
Includes returns an array of included blacklist domains/hosts



### func (Config) Sources
``` go
func (c *Config) Sources(node string) []Srcs
```
Sources returns a Leaf array for the node



### func (Config) WriteIncludes
``` go
func (c *Config) WriteIncludes(dir string, nodes []string) (dex, ex List)
```
WriteIncludes writes pre-configure data to disk



## type Configure
``` go
type Configure interface {
    Files() []string
    FormatData(node string, data []string) (reader io.Reader, list List, err error)
    Get(node string) (e *EdgeOS)
    IP(node string) string
    Sources(node string) []Srcs
    Disabled(node string) bool
    Excludes(node string) []string
    Includes(node string) []string
}
```
Configure has methods for returning config data supersets











## type Content
``` go
type Content struct {
    // contains filtered or unexported fields
}
```
Content struct holds content data









### func NewContent
``` go
func NewContent(toRead interface{}) *Content
```
NewContent returns a new Content pointer




### func (\*Content) Read
``` go
func (c *Content) Read(p []byte) (n int, err error)
```


## type Data
``` go
type Data interface {
    Reader
    Writer
}
```
Data inteface implements the Reader and Writer interfaces











## type EdgeOS
``` go
type EdgeOS struct {
    Disabled bool
    Exc      []string
    Inc      []string
    IP       string
    Nodes    []Leaf
    Sources  []Srcs
}
```
EdgeOS struct for normalizing EdgeOS data.











## type Keys
``` go
type Keys []string
```
Keys is used for sorting operations on map keys











### func (Keys) Len
``` go
func (k Keys) Len() int
```
Len returns length of Keys



### func (Keys) Less
``` go
func (k Keys) Less(i, j int) bool
```
Less returns the smallest element



### func (Keys) Swap
``` go
func (k Keys) Swap(i, j int)
```
Swap swaps elements of a key array



## type Leaf
``` go
type Leaf struct {
    Data     map[string]*Srcs `json:"data, omitempty"`
    Disabled bool             `json:"disable"`
    Excludes []string         `json:"excludes, omitempty"`
    Includes []string         `json:"includes, omitempty"`
    IP       string           `json:"ip, omitempty"`
}
```
Leaf is a struct for EdgeOS configuration data











## type List
``` go
type List map[string]int
```
List is a map of int









### func GetSubdomains
``` go
func GetSubdomains(s string) (l List)
```
GetSubdomains returns a map of subdomains


### func MergeList
``` go
func MergeList(a, b List) List
```
MergeList combines two List maps




### func (List) KeyExists
``` go
func (l List) KeyExists(s string) bool
```
KeyExists returns true if the key exists



### func (List) String
``` go
func (l List) String() string
```


### func (List) SubKeyExists
``` go
func (l List) SubKeyExists(s string) bool
```
SubKeyExists returns true if part of all of the key matches



## type Lister
``` go
type Lister interface {
    KeyExists(s string) bool
    String() string
    SubKeyExists(s string) bool
}
```
Lister implements List methods











## type Nodes
``` go
type Nodes map[string]*Leaf
```
Nodes is a map of Leaf nodes









### func NewNodes
``` go
func NewNodes() Nodes
```
NewNodes implements a new Node map


### func ReadCfg
``` go
func ReadCfg(reader io.Reader) (Nodes, error)
```
ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure




### func (Nodes) JSON
``` go
func (n Nodes) JSON() string
```
JSON returns raw print for the Blacklist struct



### func (Nodes) NewConfig
``` go
func (n Nodes) NewConfig() (c *Config)
```
NewConfig returns an initialized Config map of struct EdgeOS



### func (Nodes) SortKeys
``` go
func (n Nodes) SortKeys() (pkeys Keys)
```
SortKeys returns an array of sorted strings



### func (Nodes) SortSKeys
``` go
func (n Nodes) SortSKeys(node string) (skeys Keys)
```
SortSKeys returns an array of sorted strings



### func (Nodes) String
``` go
func (n Nodes) String() string
```
String returns pretty print for the Blacklist struct



## type Reader
``` go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```
Reader implements the Read []byte method











## type Srcs
``` go
type Srcs struct {
    Desc     string `json:"desc, omitempty"`
    Disabled bool   `json:"disabled, omitempty"`
    File     string `json:"file, omitempty"`
    IP       string `json:"ip, omitempty"`
    List     List   `json:"-"`
    Name     string `json:"name"`
    No       int    `json:"-"`
    Prefix   string `json:"prefix"`
    Type     int    `json:"type, omitempty"`
    URL      string `json:"url, omitempty"`
}
```
Srcs holds download source information









### func Process
``` go
func Process(s *Srcs, dex, ex List, reader io.Reader) *Srcs
```
Process extracts hosts/domains from downloaded raw content




## type Writer
``` go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```
Writer implements the Write []byte method

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)