<h1 align="center">rng-deviation-visualizer</h1>
<img src="https://img.shields.io/badge/go-%2300ADD8.svg?&style=for-the-badge&logo=go&logoColor=white"  align="right" position="absolute">

## Description

**rng-deviation-visualizer** was mede to visualize deviations of numbers created using a rng, the basic functionality of this project consists of generating the random numbers, sending them to the frontend using WebSocket and plotting that deviation data on a chart for visual representation.

### Project Structure

```bash
├───frontend
├───internals
│   ├───producer
│   └───websocket
└───pkg
    └───rng
```