import React from 'react'
import BarChart from '../components/BarChart'
import LineChart from '../components/LineChart'
import HeatMap from '../components/heat-map'
import GithubCorner from 'react-github-corner'
import Link from 'next/link'

export async function getServerSideProps() {

  // const base = 'http://localhost:' + process.env.PURE_PORT_BACKEND + '/people.json?group=';
  const base = 'https://gymtrackerapi.jimeagle.com/people.json?group='

  let [yearDay, monthDay, weekDay, weekHour, hour, now] = await Promise.all([
    fetch(base + 'yearDay').then((response) => response.json()),
    fetch(base + 'monthDay').then((response) => response.json()),
    fetch(base + 'weekDay').then((response) => response.json()),
    fetch(base + 'weekHour').then((response) => response.json()),
    fetch(base + 'hour').then((response) => response.json()),
    fetch(base + 'now').then((response) => response.json()),
  ])

  return {props: {yearDay, monthDay, weekDay, weekHour, hour, now}}
}

function HomePage({yearDay, monthDay, weekDay, weekHour, hour, now}) {
  return (
    <>
      <GithubCorner href="https://github.com/Jleagle/gym-tracker" bannerColor="#2f7ed8"/>
      <div className="row">
        <div className="col">

          <p>Currently recording data from Fareham only, more coming soon. <Link href="/new-gym">Add your gym</Link>.</p>

          <h2>Last 24 hours</h2>
          <LineChart data={now}/>

          <h2>By hour of the day</h2>
          <BarChart data={hour}/>

          <h2>By hour of the week</h2>
          <HeatMap data={weekHour}/>

          <h2>By day of the week</h2>
          <BarChart data={weekDay}/>

          <h2>By day of the month</h2>
          <BarChart data={monthDay}/>

          <h2>By day of the year</h2>
          <BarChart data={yearDay}/>

          <footer>If a gym has 10 or less members inside, it will show as 0.</footer>
        </div>
      </div>
    </>
  )
}

export default HomePage
