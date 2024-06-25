# Muzz 🦋

## Requirements
* Go v1.22.4 installed on your machine 
* MySQL installed on your machine 
* Docker installed on your machine if you intend to run the tests

## Setup

To set up this project, excute the setup script by running:

```bash
bash setup.sh
```

This would create the environment file prefilled with the information you provide and import the database migration.

## Running

Once setup is done, just like every simple Go project out there run the following command to start the server:

```bash 
go run server.go
```

## Usage

![server](https://res.cloudinary.com/ichtrojan/image/upload/v1719273507/server_gdzdpm.png)

once up and running, the API would be served on port `6666`

![postman](https://res.cloudinary.com/ichtrojan/image/upload/v1719273607/Screenshot_2024-06-25_at_12.59.52_am_ez9pcc.png)

You can visit the API documentation [here](https://documenter.getpostman.com/view/23911707/2sA3XY6xcG)

## Database structure

![db](https://res.cloudinary.com/ichtrojan/image/upload/v1719262721/Untitled_yvz6gf.svg)

>**NOTE**<br/>
> I have made the column names on every table self explanatory

<br/>

| Table Name | Purpose |
|---|---|
| users | This table holds user's information including location (longitude and latitude). |
| swipes | This table contains the data of users swiped on by a user. |
| swipe matches | This table holds user matches between two users, it is populated when two users mutually swipe `yes` on each other.  |


## Decision making 

### Swipes 

The Swipe implementation handles the logic when a user swipes left or right on another user in Muzz. Here's a step-by-step breakdown of what it does:

* Get the Current User: It retrieves the user making the request from the request context.

* Self-Swipe Check: It checks if the user is trying to swipe on themselves. If so, it returns an error saying, "you cannot swipe yourself".

* Check Swiped User: It looks up the user that the current user is trying to swipe on using their user ID. If the user doesn't exist, it returns an error saying, "user does not exist".

* Duplicate Swipe Check: It checks if the current user has already swiped on this user before. If they have, it returns an error saying, "you already swiped on this user".

* Handle 'No' Swipe: If the current user's swipe preference is "no" (indicating a left swipe):
   - It records the swipe in the database with a "no" preference.
   - It returns a response indicating that there is no match.
  
* Handle 'Yes' Swipe: If the current user's swipe preference is "yes" (indicating a right swipe):
  - It records the swipe in the database with a "yes" preference.
  - It then checks if the swiped user has also swiped "yes" on the current user, indicating a mutual interest.

* Mutual Interest Check:
  - If there is no mutual interest (the swiped user hasn't swiped "yes" on the current user), it returns a response indicating no match.
  - If there is mutual interest (the swiped user has also swiped "yes" on the current user), it creates a match record in the database.
  - It returns a response indicating a match, including the match ID.

I made sure that swipe implementation ensures a user cannot swipe on themselves, checks if the swiped user exists, prevents duplicate swipes, records the swipe with the appropriate preference, and determines if there is a mutual interest, creating a match if both users have swiped "yes" on each other.

### Discover

The Discover implementation provides a filtered list of potential matches for the current user, based on optional age and gender filters, and excluding users the current user has already swiped on. It also includes pagination information to help manage large sets of results.

## Testing

Execute this command to run the tests

```bash 
go test ./...
```