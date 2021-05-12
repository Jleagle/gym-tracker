import React from 'react'
import BarChart from '../components/bar-chart'
import LineChart from '../components/line-chart'
import HeatMap from '../components/heat-map'
import GithubCorner from 'react-github-corner'
import Link from 'next/link'

export async function getServerSideProps() {

  // const base = 'http://localhost:' + process.env.GYMTRACKER_PORT_BACKEND + '/people.json?group=';
  const base = 'https://gymtrackerapi.jimeagle.com/people.json?group='

  let [yearDay, monthDay, weekDay, weekHour, dayHour, now] = await Promise.all([
    fetch(base + 'yearDay').then((response) => response.json()),
    fetch(base + 'monthDay').then((response) => response.json()),
    fetch(base + 'weekDay').then((response) => response.json()),
    fetch(base + 'weekHour').then((response) => response.json()),
    fetch(base + 'dayHour').then((response) => response.json()),
    fetch(base + 'now').then((response) => response.json()),
  ])

  return {props: {yearDay, monthDay, weekDay, weekHour, dayHour, now}}
}

function HomePage({yearDay, monthDay, weekDay, weekHour, dayHour, now}) {

  return (
    <>
      <GithubCorner href="https://github.com/Jleagle/gym-tracker" bannerColor="#2f7ed8"/>
      <div className="row">
        <div className="col">

          <p>
            Currently recording data from Fareham only.&nbsp;
            <Link href="/new-gym">Add your gym</Link>.
          </p>

          <h2>Last 24 hours</h2>
          <LineChart data={now}/>

          <h2>By hour of the day</h2>
          <BarChart data={dayHour}/>

          <h2>By hour of the week</h2>
          <HeatMap data={weekHour}/>

          <h2>By day of the week</h2>
          <BarChart data={weekDay}/>

          <h2>By day of the month</h2>
          <BarChart data={monthDay}/>

          <h2>By day of the year</h2>
          <BarChart data={yearDay}/>

          <footer>Data updated every 10 minutes. If a gym has 10 or less members inside, it will show as 0.</footer>
        </div>
      </div>
    </>
  )
}

export default HomePage
