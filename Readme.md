# fdiff

**fdiff** is a tool that find the difference between files based on signatures. The tool has two commands
    - **signature** - create a signature of a file. The file is split on chunks and every chunk has its signature.
All these signatures are added to a file.
    - **delta** - find the difference between two files by using signature created from the command **signature**. 
The command returns only the different chunks of the files.

The tool uses **Rabin fingerprint** rolling hash to split the files into chunks with different size and resistant 
boarders to bytes shifting. For example if a new byte is added in the beginning of the file than all borders of the 
chunks will be shifted with one byte too and most of the chunks (except the first one) will contain the same bytes 
as before. The implementation of the Rabin fingerprint rolling hash can be found in the file [rabin_fingerprint.go](https://github.com/EmilGeorgiev/fdiff/blob/master/rollinghash/rabin_fingerprint.go)

The project contains 3 parts the first part is responsible for IO operations [signer_delta](https://github.com/EmilGeorgiev/fdiff/blob/master/signer_delta.go) 
of the files, the second one split files on chunks [chuncker](https://github.com/EmilGeorgiev/fdiff/blob/master/chuncker.go), 
and the last one contains the implementation of the rolling hash [Rabin fingerprint](https://github.com/EmilGeorgiev/fdiff/blob/master/rollinghash/rabin_fingerprint.go). 


```
    |------------------|    bytes     |------------------|             |-------------------|
    |                  | -----------> |                  |             |                   |
    | sign_delta (I/O) |   chunks     |     CHUNCKER     | ----------->| ROLLING HASH Impl.|
    |                  | <----------- |                  |             |                   |
    |------------------|              |------------------|             |-------------------|                  
```

The implementation of each part can be changes without breaking the logic in other parts (if the contract interfaces are unchanged)

## Installation

Clone the repository:
```
git clone git@github.com:EmilGeorgiev/fdiff.git
```

Install the tool:
```
go install cmd/fdiff.go
```

Now you should have **fdiff** command.

## Configuration

The file **config.yaml** contains configuration information:
- **window_size** - is the number of bytes that are included in the window that going 
to be rolling/shifted through the data.
- **min_size_chunk** - point how much must be the minimum size of a Chunk.
- **max_size_chunk** - point how much must be the maximum size of a Chunk.
- **fingerprint_break_point** - point when boundary of the chunks. 
When the hash value of the bytes in window are equal to fingerprint_break_point 
this means that the Chuncker should create a new chunk

## Example
Let's see how the tool works. First prepare a big file that you will use. For example, you can download a sample
file ("2mb text file") from here: https://www.learningcontainer.com/sample-text-file/

### Create a signature file
Run the command:
```
fdiff -signature=true -old-file sample-2mb-text-file.txt -signature-file signature
```

The command will split and sign the file **sample-2mb-text-file.txt** on chunks and a new signature file will be created. 
The command has several flags: 
- **-signature=true** - instruct the tool to create a signature file. 
- **-old-file sample-2mb-text-file.txt** - show which filed will be signed.  
- **-signature-file signature** - show in which file the signatures should be stored.

The result of the command is:
```
Creating a signature of the file:  signature
Signature file is created
```

Now you should have a file with name signature. You can see the content of the file:
```
cat signature
```

The result is:
```
0-7384-4e17f8ea25ff3a733dd03a4f8ffa68e12c7699c3
7384-27622-8fd604ec5caaa170657bc22322406fb29e3057e6
35006-10122-6e1740962a4e43c16c33d9e295306702cf8bd540
45128-25005-564ce78ffd5afb5f95182b66198eb737a76d6604
70133-17677-5770311c3af807dc2c02433cb23d57daf04e0e7d
87810-11339-7d9fd1d5b0d30f4efbdb25f76ba5e7a916326a87
99149-5393-d6b14f02eca13c7b6aaed0f67b581ff01f7e1b42
104542-10505-8669d8d7f268b1ba31e7ccffaa64941a338ee7a6
115047-21761-ee66f3c0ea0b0e50f0524a86e31fefab12fc5a51
...
```

Every line contains information about the chunks of the file **sample-2mb-text-file.txt**. Every line contains 3 parts 
separated with **-**. The first part show the offset from which the chunk started (0, 7384, 35006, ...), the second part 
shows the length of the every chunk (7384, 27622, 10122, ..) and the last part contains the signature of the chunk. 
Later in the **delta** these signatures will be used to find the differences.


### Find the delta
Now we will update the file, and we will find the difference between the two versions of the file. Let's add one new 
character **L** in the beginning of the file **sample-2mb-text-file.txt**.

```
LLorem ipsum dolor sit amet, consectetur adipiscing elit, ...
```

Now run the command with **delta** flag:
```
fdiff -delta=true -signature-file signature -new-file sample-2mb-text-file.txt
```

The command has several flags:
- **-delta=true** - instruct the tool to find the delta between two versions of the file.
- **-signature-file signature** - show which file to be used as signature. Base on this file, the tool will find the differences.
- **-new-file sample-2mb-text-file.txt** - show the new version of the file.

The result is of the command is:
```
Old chunks that are updated or removed:
	- offset: 0, length: 7384, hash: 4e17f8ea25ff3a733dd03a4f8ffa68e12c7699c3
New chunks that replace the old ones:
	- offset: 0, length: 7385, hash: 9b8d14a2408f987a136c6c414b7aea1ddc4b7238
```

You can see the old chunks that are updated or removed and are not up-to-date. You can see and the new chunks that will 
replace the old.

Also, if you want you can see the data in the new chunks by using the flag **show-data**
```
fdiff -delta=true -signature-file signature -new-file sample-2mb-text-file.txt -show-data=true
```

The result is:
```
Old chunks that are updated or removed:
	- offset: 0, length: 7384, hash: 4e17f8ea25ff3a733dd03a4f8ffa68e12c7699c3
New chunks that replace the old ones:
	- offset: 0, length: 7385, hash: 9b8d14a2408f987a136c6c414b7aea1ddc4b7238
	- LLorem ipsum dolor sit amet, consectetur adipiscing
	...
```
