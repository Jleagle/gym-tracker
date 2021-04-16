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
                    <h1 className="mt-4">PureGym Tracker</h1>
                    {/*<LineChart data={now}/>*/}
                    {/*<BarChart data={hour}/>*/}
                    {/*<HeatMap data={weekHour}/>*/}
                    {/*<BarChart data={weekDay}/>*/}
                    {/*<BarChart data={monthDay}/>*/}
                    <BarChart data={yearDay}/>
                </div>
            </div>
        </div>
    );
}

export default HomePage
