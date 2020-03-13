# vmwarego

Small wrapper around [govmomi](github.com/vmware/) for getting VM info simpler

Connections parameters and client:
```go
vmw, err := NewClient(ClientParams{
    URL:      "https://HOST",
    Insecure: true,
    User:     "USER",
    Password: "PASS",
    Ctx:      context.Background(),
})
```

#### Some examples

Find VM by name:
```go
vm, err := vmw.VmInfo("golden-w2k16")
// vm =>
/*
{
    "Vm": {
      "Type": "VirtualMachine",
      "Value": "vm-152148"
    },
    "Runtime": {
      "Device": [
        {
...
*/
```

Find VMs by a mask and get a custom field:
```go
vms, err := vmw.VmsFilter("base-w2k16-*")
if err != nil {
    t.Fatal(err)
}

err = vmw.VmRetrieve(vms, []string{"summary"})
if err != nil {
    t.Fatal(err)
}

type Tags struct {
    Template string `json:"Template"`
}
var tags Tags

for _, vm := range vms {
    err = vmC.VmCustomFields(vm, &tags)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Printf("%s || Template => %s\n", vm.Summary.Config.Name, tags.Template)
}
// vm =>
/*
base-w2k16-123 || Template => true
base-w2k16-1234 || Template => true
base-w2k16-12345 || Template => true
...
*/
}
```

