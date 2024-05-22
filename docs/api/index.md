# Introduction
**Authors:** Scott Barnard

OpenFish provides an API to access stored marine footage, and video annotations / labels, allowing clients to retrieve and filter the data.
Clients can download segments of footage or video annotations by querying by location, time, and other parameters.

OpenFish's API has a few types of resources it deals with: [capture sources][capture-sources], [video streams][video-streams], [annotations][annotations], [species][species] and [users][users]. Capture sources are cameras that produces video streams. Video streams have information about a single video. Annotations are used for labeling interesting things at a particular time and place in videos. Species provides a list of valid species for users to identify. Users contains the user's [role and permissions][roles-permissions].

[capture-sources]: ./capture-sources
[video-streams]: ./video-streams
[annotations]: ./annotations
[species]: ./species
[users]: ./users
[roles-permissions]: ./roles-and-permissions

