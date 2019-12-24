package database

const (
	createRoutesTable = `
            CREATE TABLE %s(
                id               INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                route_id         VARCHAR(10) NOT NULL,
                route_short_name VARCHAR(5) NOT NULL,
                route_long_name  VARCHAR(255) NOT NULL,
                route_desc       VARCHAR(255) NOT NULL
            ) ENGINE=InnoDB
        `

	createStopsTable = `
            CREATE TABLE %s(
                id              INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                stop_id         INT NOT NULL,
                stop_code       INT NOT NULL,
                stop_name       VARCHAR(55) NOT NULL,
                stop_desc       VARCHAR(55) NOT NULL
            ) ENGINE=InnoDB
        `

	createStopTimesTable = `
            CREATE TABLE %s(
                trip_id         INT NOT NULL,
                arrival_time    VARCHAR(9) NOT NULL,
                departure_time  VARCHAR(9) NOT NULL,
                stop_id         INT NOT NULL
            ) ENGINE=InnoDB
        `

	createTripsTable = `
            CREATE TABLE %s(
                route_id     VARCHAR(10) NOT NULL,
                service_id   VARCHAR(30) NOT NULL,
                trip_id      INT NOT NULL,
                direction_id INT NOT NULL
            ) ENGINE=InnoDB
        `

	createCalendarTable = `
            CREATE TABLE %s(
                service_id VARCHAR(30) NOT NULL,
                monday     VARCHAR(15) NOT NULL,
                tuesday    VARCHAR(15) NOT NULL,
                wednesday  VARCHAR(15) NOT NULL,
                thursday   VARCHAR(15) NOT NULL,
                friday     VARCHAR(15) NOT NULL,
                saturday   VARCHAR(15) NOT NULL,
                sunday     VARCHAR(15) NOT NULL,
                start_date VARCHAR(10) NOT NULL,
                end_date   VARCHAR(10) NOT NULL
            ) ENGINE=InnoDB
        `
)
