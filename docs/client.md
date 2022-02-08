# Store Client

The store client is a command line utility that works with a store server endpoint to perfomr the operations on the store.

## Commands

### **config**
Config command will ask for the remote store URL and save it in the local filesystem as .StoreConfig file. A valid store URL is required else the command will fail. This also automatically executed if the before executing any other command if the .StoreConfig file is missing

```
$ ./store config
Store URL: http://localhost:4000
```

### **add** 
Add files in the remote store. Existing file names cannot be added 
  
```
$ ./store add file1.txt file2.txt

file1.txt - Added successfully
file2.txt - Failed to upload
```

### **update** 
Update existing file in the remote store. If a file does not exists, it will be created.  

```
$ ./store update file2.txt file3.txt

file2.txt - Updated successfully
file3.txt - Added successfully
```

### **remove** 
Remove a file from the remote store
  
```
$ ./store remove file2.txt

File deleted successfully!
```

### **list** 
List the files on the remote store. The command supports an extra flag **--details**, which shows more information on the file.

```
$ ./store list --details

FILE NAME                 BYTES      WORDS
---------                 -----      -----
file1.txt                   26          5
file3.txt                   26          5
```

### **get**
Get command will fetch a file from the server and save in the local filesystem

```
$ ./store get file1.txt
```
### **count**
This command display the total number of words in all the files on the remote store.

Can also be accessed with the alias **wc**

```
$ ./store wc
Word Count: 64
```
### **frequency**
This command will display the frequency of all the words in the all the files combined on the remote store.

This takes the the below optional flags..
- --limit : will limit the number of results
- --order : will return ascending or descending order

```
$ ./store freq-words --limit 5 --order dsc
        12 file
         6 for
         6 test
         6 sample
         6 date
```

## List of all commands

```
Usage:
  store [command]

Available Commands:
  add         Add files
  completion  Generate the autocompletion script for the specified shell
  config      Configure store
  count       Word count
  frequency   Word frequency
  get         Get File
  help        Help about any command
  list        List files
  remove      Remove files
  update      Update/Create files
  version     Store version

Flags:
  -h, --help      help for store
  -v, --version   version for store

Use "store [command] --help" for more information about a command.
```