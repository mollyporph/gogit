# GoGit

GoGit is a small github v3 api client specifically built to ease overview for organisations with multiple repos.
![prpersonal](https://puu.sh/s8oRm/7991789e8e.png)

#Features
* Using a personal access token instead of requiring username and password for each call
* Works with two-factor authentication
* Get pullrequests for organisation, team or personal and filter them by state
* Pullrequests can be sorted by all fields


#Todo
* Get issues for organisation, team or personal
* Handle paging correctly from github
* Colorize the table (for instance red for pullrequests that has been open for more than x days)
* Add config fields for default warning threshold for pullrequest open duration
* Add ability to modify config file after initial setup
* Add ID to pullrequest table
* Add detailview for a single pullrequest (i.e gogit pr 12345)
* Add detailview for a single issue (i.e gogit is 12345)

