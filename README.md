# KIN

## Architecture
So the project essentially sends system information to a theoretically infinite amount of keyboards.
To handle this we should fetch information once per refresh rate, and then the same information to each 
keyboard. From a configuration perspective, this should look something like this:

Each keyboard we want to send information to will have its own information which it wants. To resolve this
our config file should store each keyboard with its own internal refresh rate, which is how often we send information
that specific keyboard. On top of this, each piece of information we fetch should have its own separate refresh rate, which
defines how often we retrieve fresh data from the system

To this extent, our config should be split into 2 different major parts. The information fetching, and the 
information sending.

Each piece of information we fetch should have the ability to be blacklisted from fetching for those who are privacy minded,
 and on top of that each keyboard defines what information it wants.

For keyboards defining information which it wants, we could possibly create an init system to ask the keyboard
which pieces of information it wants, and then let the user decide which pieces of information it gets. Realistically this
process shouldn't be completely automatic since some information may be considered sensitive.