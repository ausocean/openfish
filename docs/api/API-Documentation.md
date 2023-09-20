# Overview
OpenFish provides an API to access stored marine footage, and video annotations / labels, allowing clients to retrieve and filter the data.
Clients can download segments of footage or video annotations by querying by location, time, and other parameters.

OpenFish's API has three types of resources it deals with: [capture sources][capture-sources], [video streams][video-streams] and [annotations][annotations]. Capture sources are cameras that produces video streams, video streams have information about a single video and annotations are used for labeling interesting things at a particular time and place in videos. 

The OpenFish API uses a REST pattern and has common patterns for certain tasks for all types of resources. See [General Usage Notes][general-usage-notes] for more information.

[general-usage-notes]: https://github.com/ausocean/openfish/wiki/general-usage-notes
[capture-sources]: https://github.com/ausocean/openfish/wiki/capture-sources
[video-streams]: https://github.com/ausocean/openfish/wiki/video-streams
[annotations]: https://github.com/ausocean/openfish/wiki/annotations

