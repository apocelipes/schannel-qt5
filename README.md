## schannel-qt5
A QT based GUI client for [schannel](https://schannel.net/) - written in Golang

### Features:
- clear data usage display
- quickly view service information
- supports lots of user settings
- simply & easy to use

### Installation
At first we need to install [thetheripe/qt](https://github.com/therecipe/qt)

Then:
```bash
cd $GOPATH/src
git clone 'https://github.com/apocelipes/schannel-qt5'
cd schannel-qt5/widgets && qtmoc
cd .. && go build
```

Now you can enjoy schannel-qt5!

### Screenshots
login:

![login](screenshots/login.png)

invoices view:

![invoices](screenshots/invoices.png)

select nodes:

![nodes](screenshots/nodes.png)

service info & client switch:

![service](screenshots/nodes.png)

user settings:

![settings](screenshots/settings.png)

### Todo:
- add charts to display daily usage
- add delete user button in LoginWidget

Welcome feedback questions and submit PRs,

I am looking forward to working with you.