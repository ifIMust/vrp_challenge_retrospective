package common

// HomeLocation is the defined as depot location.
// All drivers must start and end their route here.
var HomeLocation *Location = &Location{0, 0, 0.0}

// MaxMinutesPerDriver is the total travel time per day allowed for any one driver.
const MaxMinutesPerDriver = 12 * 60
