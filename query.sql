        SELECT DISTINCT
            s.stop_name, -- from
            s2.stop_name, -- to
            st.arrival_time, -- from time
            s2.arrival_time, -- to time
            r.route_short_name
        FROM stops s
        INNER JOIN stop_times st
            ON st.stop_id = s.stop_id
        INNER JOIN trips t
            ON t.trip_id = st.trip_id
        INNER JOIN routes r
            ON r.route_id = t.route_id
		INNER JOIN calendar c
			ON c.service_id = t.service_id
        INNER JOIN (
            SELECT DISTINCT
                s.stop_name,
                st.arrival_time,
                t.trip_id,
                r.route_short_name,
                r.route_id
            FROM stops s
            INNER JOIN stop_times st ON st.stop_id = s.stop_id
            INNER JOIN trips t ON t.trip_id = st.trip_id
            INNER JOIN routes r ON r.route_id = t.route_id
			INNER JOIN calendar c ON c.service_id = t.service_id
        ) s2
            ON r.route_id = s2.route_id
            AND s.stop_name != s2.stop_name
            AND t.trip_id = s2.trip_id

        WHERE
            s.stop_name = 'Alameda Station'
            AND s2.stop_name = '16th & California Station'
            AND st.arrival_time < s2.arrival_time -- From time < To time
            AND st.arrival_time > 19:04:00
			AND c.Wednesday = 1
        ORDER BY st.arrival_time ASC
        LIMIT 5
