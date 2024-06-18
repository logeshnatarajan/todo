# Uncrew Architecture
This is the repository where the Uncrew architecture lives. As dictated by our [architectural process](https://github.com/droneup/common-architecture/tree/trunk/ADRs/COMMON-1), architecture is **manifested** by a set of [Architectural Decision Records (ADRs)](./ADRs/). This README is simply their directory, glue and welcome.

Uncrew is a standalone product responsible for planning and flying UAV (drone) missions. It encompasses the entire flight control vertical, initially connecting the operator (pilot) with the UAV's flight controller firmware (e.g.: [PX4](https://px4.io/)) - later gradually replacing the pilot and climbing the 5 levels of [UAV autonomy](https://dronelife.com/2019/03/11/droneii-tech-talk-unraveling-5-levels-of-drone-autonomy/).

## Terminology

**Operator** Anyone authorized to directly influence the aircraft's execution of any procedure (e.g. flight path, delivery, survey plan, etc.) during the mission. Missions execute automatically (not yet autonomously) and it's the Operator's job to oversee them and take them out of trouble. An Operator never flies with sticks - the majority of flight time any given operator is not close enough to judge the pose of the UAV. Instead; Operator flies with what ARDU Pilot calls the [guided flight mode](https://ardupilot.org/plane/docs/guided-mode.html), i.e.: stop, fly there, land. The majority of the time, there will be a single operator overseeing a given mission.  However, some circumstances (environmental, visibility or complexity) may require another or two operators to be in/on the loop with that initial operator forming the **Mission Team**.

**Supervisor** reviews the Mission plan for santity, confirms it is good to go and assigns an operator. Supervisor may, in some circumstances, reassign a new operator.

**Apollo** is the cloud part of Uncrew named so because it's essentially the mission control.

**Site** (a.k.a. **Hub**) - Is a mission dispatch centre (i.e.: location rather than a region). This is where a number of UAVs are charged, loaded with delivery packages and then flown by a group of operators sitting in front of their control stations. UAVs can migrate from one hub to another with an easy, manual or automated (drone-location-based) procedure, but at any given time, a UAV is associated with one hub. This constraint isn’t relevant for delivery missions where it’s the HubOps that appoint the UAV to fly. It becomes relevant when a client relies on Apollo to pick a UAV for execution and Apollo needs start doing fleet management.

**Operational Area** A defined area within a boundary where Missions occur for an operations group within an organization.  Sites, Launch Zone points, Land Zone points and other pre-defined structures live within this boundary. The "boundary" defines the area where the UAS Implementation team has scouted and performed risk assessments etc. within the geography. Encompasses the ground infrastructure assets (logistics centers AKA hubs, destination boxes, etc) and their corresponding LZ (Launch/Land Zone) points. Mission Teams and Vehicles are assigned within an Operational Area. If a Mission comes into a logistics center (AKA a hub) that asset is within an Operational Area, and the mission is generated within that subset of geodata.

## Apollo Mission Control
Before the operator can be taken away, it needs to be first put in. Apollo is the first seed of Uncrew - a mission-control application that connects (human) operators, missions and UAVs to enable the DroneUp drone services (including delivery!) use case.
![alt text](./doc/uncrew-architecture.drawio.png)

Modular by separating the UI from the backend and further by modularizing the backend at the service level. Ultimately, to satisfy the DroneUp Delivery use-case, Apollo Mission Control should be deployed on a site local hardware to minimize network latency and mitigate the risk of upstream network outage. As a bonus, Apollo, with a subset of services, should be also deployable to a single device (e.g.: an Linux/Windows tablet, though Android would be [a stretch](https://stackoverflow.com/questions/53527277/is-it-possible-to-run-containers-on-android-devices) and iOS a no go) becoming a traditional, self-contained ground control station (we will refer to it as the self-contained GCP use case). Meanwhile it is simply a set of microservices deployed in a Kubernetes cluster.

### Frontend
Because the Visual Observer is mobile and Ansur encodes the drone’s video feed with an exotic codec, at least one of the Apollo frontends should be native, yielding difficult, product-level choices: Linux (e.g.: GtK), Windows, Android or iOS? Native, React Native or Qt? Meanwhile the accessibility argument dictates that the one UI Apollo should have (probably the first one to build) is **web** - your experience will be a bit worse, but you can fly the drones on any platform.

### UTM
Uncrew depends on UTM in the following areas.
#### Geospatial Data
Most (if not all) data that the pathinder needs to pay attention to (when plotting paths) comes from a/the UTM: terrain, airspaces, previously submitted flight volumes and (at some point in time hopefully) population density. The data is to be portioned into tiles akin to mapbox tiles. Uncrew will set up a local and on-UAV cascading cache to be in control of the data avaialbility.

#### Flight Submissions
For every mission, Uncrew will submit a flight intent to UTM. UTM shall validate the plan, check for conflicts, reserve the airspace and notify the public about the imminent flight. Uncrew will only execute the intent if UTM accepts it. Uncrew will then stream positional telemetry to UTM for publishig and conformance monitoring.
