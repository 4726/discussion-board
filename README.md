# discussion-board

Implementation of a discussion board using microservice architecture. 

Backend written in Go and frontend written in Typescript using Angular.

# Architecture

![img](https://github.com/4726/discussion-board/blob/master/images/architecture.png?raw=true)

# Services

### Frontend

Client requests are sent to the frontend and the frontend communicates with the API Gateway.

### API Gateway

Communicates with backend services using gRPC and communicates with frontend using REST.

### Likes

Stores amount of likes each post and comment has.

### Media

Stores user uploaded images.

### Posts

Stores user posts and comments.

### Search

Search for posts containing specific text.

### User

Stores user credentials and profile information.