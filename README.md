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

The technical side of the project is built various languages and frameworks:

* Node.js and Go (Gin) to develop server and bots
* ReactJS and Typescript for making user interface
* AWS (EC2) for server deployment



<p align="right">(<a href="#top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Usage

Here we describe how to use the botnet.

### To use the botnet


1. Install Go if necessary
2. You can press **"Start testing"** and the server would be active for new bots connections.
3. To run bots, run process manager (```main.go``` in ```bots/processManagerCloud```)
4. Once the bot is connected, you will see the stats on the dashboard updated every 10 seconds.
5. After all the work is finished, you can press "Stop" to deactivate the server and "Exit" to log out of the dashboard.


### To use the botnet as a developer

_Below is an example of how you can download and change the source code._

1. Clone the repo
   ```sh
   git clone https://github.com/kkulykk/distributed-botnet.git
   ```
2. Install Node.js and Go if needed and open the cloned project directory in IDE. 
3. Go to the desired folder and file you want to work with and observe the source code of the project.

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap


- [x] Design the structure and define the logic of the server and bots
- [x] Create server and bot (MVP) for debugging with Node.js
- [x] Master Go to implement both bots and server
- [x] Create UI for managing the botnet
- [x] Establish SSH connection on local machines for the further bots updating
- [x] Deploy server on AWS for testing and demonstration purposes
- [x] First presentation
- [ ] ~~Adding fuctionality to test multiple targets~~
- [x] Add `Requests Amount` and `Time attack` mode
- [x] Remote bots updating using ~~SSH~~ Amazon S3
- [x] Implementing bots in ~~C++~~ Go using goroutines
- [x] Debugging and testing
- [x] Final presentation


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


