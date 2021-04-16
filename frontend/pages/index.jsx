import React from 'react'
import BarChart from "../components/BarChart";

export async function getServerSideProps() {

    const base = 'https://pgt2.jimeagle.com/people.json?group=';

    let [yearDay, monthDay, weekDay, weekHour, hour, now] = await Promise.all([
        fetch(base + 'yearDay').then(response => response.json()),
        // fetch(base + 'monthDay').then(response => response.json()),
        // fetch(base + 'weekDay').then(response => response.json()),
        // fetch(base + 'weekHour').then(response => response.json()),
        // fetch(base + 'hour').then(response => response.json()),
        // fetch(base + 'now').then(response => response.json()),
    ]);

    return {props: {yearDay}};
}

function HomePage({yearDay, monthDay, weekDay, weekHour, hour, now}) {

    return (
        <div className="container">
            <div className="row">
                <div className="col">
                    {/*<h1 className="mt-4">PureGym Tracker</h1>*/}

                    <h2>Last 24 hours</h2>
                    {/*<LineChart data={now}/>*/}

                    <h2>By hour</h2>
                    {/*<BarChart data={hour}/>*/}
                    {/*<HeatMap data={weekHour}/>*/}

                    <h2>By day of the week</h2>
                    {/*<BarChart data={weekDay}/>*/}

                    <h2>By day of the month</h2>
                    {/*<BarChart data={monthDay}/>*/}

                    <h2>By day of the year</h2>
                    <BarChart data={yearDay}/>

                    <footer>Footer</footer>
                </div>
            </div>
        </div>
    );
}

export default HomePage
