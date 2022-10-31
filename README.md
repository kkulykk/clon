# Clon

<div id="top"></div>


<br />
<div align="center">
    <img src="https://img.freepik.com/free-vector/cute-elephant-sitting-waving-hand-cartoon-vector-icon-illustration_138676-2220.jpg?w=2000" alt="Logo" width="200" height="auto">

  <h3 align="center">Clon (an <a href="https://rclone.org/">Rclone</a> analogue)</h3>

  <p align="center">
    Third-term <b>Operational Systems</b> course project 
    <br />
        <br />
    <a href="https://github.com/kkulykk/clon/issues">Request Feature</a>
  </p>
    <br />
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contacts">Contacts</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project


The project is an attempt to create an Rclone analogue. Rclone is a command-line program to manage files on cloud storage. It is a feature-rich alternative to cloud vendors' web storage interfaces. Over 40 cloud storage products support rclone including S3 object stores, business & consumer file storage services, as well as standard transfer protocols. Rclone has powerful cloud equivalents to the unix commands `rsync`, `cp`, `mv`, `mount`, `ls`, `ncdu`, `tree`, `rm`, and `cat`. It is used at the command line, in scripts or via its API.

In our implementation, called Clon (pronounced s-lon) we aimed to focus on intuitive command line interface with less configurable settings. As a result, we built a command-line tool to easily manage data on cloud storage and transfer files between local and remote storage instances.


<p align="right">(<a href="#top">back to top</a>)</p>



### Built With

The technical side of the project is built using `GoLang` only.


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- GETTING STARTED -->
## Usage

Here we describe how to use the program.

### To use `Clon`


1. Install Go if necessary
2. Set up your AWS S3 credentials in `.env` in order to establish connection with remote
3. Install the required dependincies
4. Build the project and run one of the available commands

### Available commands

#### about
Usage: `clon [OPTIONS] about`

Get detailed information about the Clon program.

![Screen Recording 2022-10-31 at 08 30 06](https://user-images.githubusercontent.com/72144618/198945553-d36147af-554d-4215-8cf9-9cf3e288062d.gif)

#### create-remote
Usage: `clon [OPTIONS] create-remote RemoteName`

Create a new remote with given name

![Screen Recording 2022-10-31 at 08 55 15](https://user-images.githubusercontent.com/72144618/198949280-7c30e363-53b3-446e-87c1-59ed83303882.gif)

#### delete-remote

Usage: `clon [OPTIONS] delete-remote RemoteName`

Delete an existing remote.

![Screen Recording 2022-10-31 at 08 52 44](https://user-images.githubusercontent.com/72144618/198948965-82d0fd78-feb1-4141-acbd-8fe1d8762b13.gif)

#### copy

Usage: `clon [OPTIONS] copy [copy-OPTIONS] FromPath ToPath`

Copy file(s) or directories from remote or local or vice versa. Use -f to force
copy.

![Screen Recording 2022-10-31 at 19 55 26](https://user-images.githubusercontent.com/72144618/199076537-47e55511-1c24-400b-88d6-3449d125645e.gif)

#### move

Usage: `clon [OPTIONS] move FromPath ToPath`

Moves the contents of the source directory to the destination directory

![Screen Recording 2022-10-31 at 19 57 39](https://user-images.githubusercontent.com/72144618/199077076-a2a83ccd-7270-4da1-a3bc-5f1edfbdc1de.gif)

#### delete

Usage: `clon [OPTIONS] delete RemotePath`

Lists the objects in the source path to standard output in a human readable
format with size path and update date.

![Screen Recording 2022-10-31 at 10 00 49](https://user-images.githubusercontent.com/72144618/198960291-0076a915-be17-4fa9-aa86-6d2e165152f3.gif)

#### ls

Usage: `clon [OPTIONS] ls Path`

Remove the files in path.

![Screen Recording 2022-10-31 at 08 58 21](https://user-images.githubusercontent.com/72144618/198949662-a52a4db7-374f-41b4-8322-94bc450b6be4.gif)


#### size

Usage: `clon [OPTIONS] size FilePath`

Return size of file in bytes.

![Screen Recording 2022-10-31 at 20 02 35](https://user-images.githubusercontent.com/72144618/199077871-267e9c93-62d8-4e49-b20c-8a1e177326db.gif)


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap


- [x] Create a simple go project
- [x] Add command line options parser
- [x] Set up access with Amazon S3 API
- [x] Add `create-remote` command
- [x] Add `delete-remote` command
- [x] Add `copy command` support
- [x] Add `move` command support
- [x] Add `delete` command support
- [ ] ~~Add `mkdir` command support~~ (impossible as for now)
- [x] Add `ls` command support 
- [x] Add command for listing all remotes
- [x] Add `size` command support
- [ ] ~~Add `delete-file` command support~~ (was implemented inside `delete` command)
- [x] Add `about` command support
- [x] First presentation
- [ ] Add `sync` command support
- [ ] Add `check` command support
- [ ] Add multiple providers support
- [ ] Add SHA-1 encryption


<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTACT -->
## Contacts

Bohdan Mykhailiv - [Github](https://github.com/bmykhaylivvv)

Roman Kulyk - [Github](https://github.com/kkulykk)


<p align="right">(<a href="#top">back to top</a>)</p>


