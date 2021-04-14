import React from 'react'
import HeatMap from '../components/HeatMap.jsx'
import LineChart from '../components/LineChart.jsx'
import BarChart from '../components/BarChart.jsx'

export async function getServerSideProps() {

    const res = await fetch(`http://localhost:` + process.env.PURE_PORT_BACKEND + `/heatmap.json`)
    const data = await res.json()

    return {props: {data}}
}

function HomePage({data}) {

    return (
        <div className="container">
            <div className="row">
                <div className="col">
                    <h1 className="mt-4">PureGym Tracker</h1>
                    <HeatMap data={data}/>
                    <LineChart/>
                    <BarChart/>
                </div>
            </div>
        </div>
    );
}

export default HomePage
