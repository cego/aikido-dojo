# Smoke testing

AutoShip performs a smoke test on a live environment to verify the changes

AutoShip validates the changes are not impacting the functionality of the application by running a smoke test on a live environment.

Smoke tests are configured per respository and environment, for each environment you specify:

* the domain AutoShip should connect to
* the credentials AutoShip should use to authenticate
* the instructions to follow to execute the smoke test

AutoShip will tell you wether the smoke test passed or failed and give you a report with screenshots on what it actually tested.
