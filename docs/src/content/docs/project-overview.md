---
title: Project Overview
---

**Authors:** Trek Hopton, Alan Noble


## Mission
AusOcean is dedicated to helping our oceans through the use of technology. Our development of a video livestream stack has resulted in a significant accumulation of marine footage being uploaded to YouTube, with hours of content available for viewing. While the availability of this marine content contributes to the raising of awareness for our oceans, we believe the video data provides significant potential for further positive impact.

The goal of the OpenFish project is to improve our understanding of marine species.

## Description
OpenFish is an open-source system written in Golang for classifying marine species. Tasks involve importing video or image data, classifying and annotating data (both manually and automatically), searching, and more. It is expected that OpenFish will utilise computer vision and machine learning techniques.

## Stages

### Step 1 - AusOcean Data Storage Method
AusOcean will need to provide access to its marine video content to start the project off. Currently AusOcean livestreams to YouTube, on which the data is stored on the AusOcean channel. The data needs to be in the google datastore in preparation for a data access API. 

Initially a collection of videos should be downloaded in mp4 (h264) from YouTube. This can happen via the YouTube Studio UI or using existing libraries such as pytube. Then we convert to MPEG-TS using ffmpeg, then upload to the datastore using VidGrind’s upload page. This could be automated but it’s not super important as this is just a starting point.

VidForward is AusOcean’s cloud-hosted video stream middleman. This will give us the ability to stream to YouTube while simultaneously streaming MPEG-TS (h264) to the datastore. Once this is implemented, we won’t have to upload videos from YouTube.

### Step 2 - Data Access API
In this step, an API will be developed to provide access to the stored marine footage, allowing clients to retrieve and use the data.

A client should be able to find and download a segment of footage by providing a query eg. location, date, time.

It would be useful to be able to extract an image from a given video, this would require h264 decoding.

### Step 3 - UI for and Video Annotation and Labelling
In this step, a user interface will be created to help users identify different species of marine life found in the footage. 

This will incorporate a video player and a labelling interface. Users should be able to set the time range for an observation and draw a box around the object of interest. It would be helpful to have a chart/guide which informs species classification.

It is desirable for users to be able to share clips of video.

If a user is not logged in, they should only be able to play the video and see the labels (read only). When logged in, users can label the data. It would be useful to have different user classes eg. certified experts have higher weights on their classifications.

### Step 4 - Dataset Curation
In this step, data from the labelling UI will be curated and organised into a training dataset for computer vision techniques.

### Step 5 - Motion Detection 
In this step, computer vision techniques will be applied to the marine footage to detect movements in the videos. Motion detection will help us narrow down the video to moments of interest.

### Step 6 - Object Detection
In this step, computer vision techniques will be applied to the marine footage to detect objects eg. fish in the videos. Object detection will allow us to automatically create bounds around marine species for more efficient manual and automatic classification.

### Step 7 - Automatic Species Classification / Suggestion
In this step, artificial intelligence algorithms will be used to automatically classify marine species based on the data collected in previous steps. A system will also be developed to provide suggestions in the UI for species classification to help users identify different species of marine life.

### Step 8 - Species Statistics and Analysis Tools
In this step, statistics and analysis tools will be developed to enable the analysis of the data collected on marine species, and to provide trends and other insights into the distribution and diversity of species in the ocean. 


## Related Projects
**Classification using images**

FishNet and FishID are both open-source frameworks for classifying fish species from images. FishNet is a Python package where FishID is a mobile phone app (iOS and Android).

(iNaturalist)[https://www.inaturalist.org/pages/developers] provides image datasets for the purpose of species classification.

**Classification using video**

(Fish4Knowledge)[https://homepages.inf.ed.ac.uk/rbf/fish4knowledge/overview.htm], a system for fish recognition and behaviour analysis from underwater video written in Python. It was a research project at the University of Edinburgh and does not appear to have been active for ten years. The code can be found on SourceForge (here)[https://sourceforge.net/projects/fish4knowledgesourcecode/]. The datasets can be found on GitHub but not the code. It uses a GNU General Public License 2.0.

It appears to use older computer-vision techniques and is designed to run on desktop computers (not the cloud), and uses the MySQL database for storage.
Useful projects

**Useful Projects**

(DeepFish)[https://github.com/alzayats/DeepFish] is described as “A Realistic Fish-Habitat Dataset to Evaluate Algorithms for Underwater Visual Analysis”. It uses an MIT licence.


## Licensing

While the GNU General Public License version 2.0 (GPLv2) and the BSD 3-Clause License are both widely-used open-source software licences, the latter will be used for OpenFish.

Key differences between the two licences are as follows:

1. Copyleft vs. Permissive: The GPLv2 is a copyleft licence, which means that any modifications or derivative works of the software must also be distributed under the same licence. This ensures that the software remains free and open-source. The BSD 3-Clause License, on the other hand, is a permissive licence that allows for modifications and derivative works to be distributed under different licences, including proprietary licences.
2. Patent Rights: The GPLv2 includes a patent clause that requires anyone who distributes the software to grant all recipients a licence to any patents that are necessary to use the software. The BSD 3-Clause License does not include a patent clause.
3. Attribution: The BSD 3-Clause License requires that any distribution of the software include a notice that the software is licensed under the BSD 3-Clause License and include the copyright notice and disclaimer. The GPLv2 also requires that any distribution of the software include a notice that the software is licensed under the GPL, but it also requires that the source code be made available to recipients and that any modifications or derivative works be clearly marked as such.
4. Compatibility: The GPLv2 and the BSD 3-Clause License are both compatible with other open-source licences, but the copyleft nature of the GPLv2 may make it less compatible with proprietary licences.
