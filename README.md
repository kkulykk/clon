# Clon

<div id="top"></div>


<br />
<div align="center">
    <img src="https://img.freepik.com/free-vector/cute-elephant-sitting-waving-hand-cartoon-vector-icon-illustration_138676-2220.jpg?w=2000" alt="Logo" width="200" height="auto">

  <h3 align="center">Clon (an <a href="https://rclone.org/">RClone</a> analogue)</h3>

  <p align="center">
    Third-term <b>Operational Systems</b> course project 
    <br />
        <br />
    <a href="https://kkulykk.github.io/distributed-botnet/" target="_blank">Go to website</a>
    ·
    <a href="https://youtu.be/ORd-A4XrvpA">View Demo</a>
    ·
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


A botnet is a collection of internet-connected devices, including personal computers (PCs), servers, mobile devices, and internet of things (IoT) devices infected and controlled by a common type of malware, often unbeknownst to their owner. However, the last part is not about our project, as it was created only for non-harmful educational purposes.

The main goal of our project is to send requests to given targets so that we can see how well the targets handle such a minor DDoS attack and how efficient the performance of network services is.

The core principle of work is also pretty straightforward: there is one central server and lots of bots (PCs or other devices) that might establish a connection with the server. The server performs some logic on gathering information on the target, how many bots are connected, what they have to do, and how to distribute the work between bots efficiently. Bots, on their side, get from the server the algorithm of what they have to do and perform those actions, periodically sending stats to the server. All the information is gathered, processed, and displayed on UI so that users can follow how the process is going on.


<p align="right">(<a href="#top">back to top</a>)</p>



### Built With

The technical side of the project is built using `GoLang` only.


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- GETTING STARTED -->
## Usage

Here we describe how to use the botnet.

### To use `Clon`


1. Install Go if necessary
2. Set up your AWS S3 credentials in `.env` in order to establish connection with remote
3. Build the project and run one of the available commands

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

Delete an existing remote

![Screen Recording 2022-10-31 at 08 52 44](https://user-images.githubusercontent.com/72144618/198948965-82d0fd78-feb1-4141-acbd-8fe1d8762b13.gif)


#### copy


#### move


#### delete

#### ls
Usage: `clon [OPTIONS] ls Path`

Lists the objects in the source path to standard output in a human readable
format with size path and update date.

![Screen Recording 2022-10-31 at 08 58 21](https://user-images.githubusercontent.com/72144618/198949662-a52a4db7-374f-41b4-8322-94bc450b6be4.gif)


#### size

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


