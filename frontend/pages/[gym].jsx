import React, {useEffect} from 'react'
import BarChart from '../components/bar-chart'
import LineChart from '../components/line-chart'
import HeatMap from '../components/heat-map'
import GithubCorner from 'react-github-corner'
import Gyms from "../components/gyms"
import {useRouter} from 'next/router'

export async function getServerSideProps() {

  // const base = 'http://localhost:' + process.env.GYMTRACKER_PORT_BACKEND + '/'
  const base = 'https://gymtrackerapi.jimeagle.com/'

  let [gyms, now, dayHour, weekHour, weekDay, monthDay, yearMonth] = await Promise.all([
    fetch(base + 'gyms.json').then((response) => response.json()),
    fetch(base + 'people.json?group=now').then((response) => response.json()),
    fetch(base + 'people.json?group=dayHour').then((response) => response.json()),
    fetch(base + 'people.json?group=weekHour').then((response) => response.json()),
    fetch(base + 'people.json?group=weekDay').then((response) => response.json()),
    fetch(base + 'people.json?group=monthDay').then((response) => response.json()),
    fetch(base + 'people.json?group=yearMonth').then((response) => response.json()),
  ])

  return {props: {gyms, now, dayHour, weekHour, weekDay, monthDay, yearMonth}}
}

function HomePage({gyms, now, dayHour, weekHour, weekDay, monthDay, yearMonth}) {

  const router = useRouter()

  useEffect(() => {

    const {gym} = router.query

    const allButtons = document.querySelectorAll('#gyms button')
    for (let i = 0; i < allButtons.length; i++) {
      const enabled = (allButtons[i].getAttribute('data-gym') === gym)
      allButtons[i].classList.add(enabled ? 'btn-success' : 'btn-link')
      allButtons[i].classList.remove(enabled ? 'btn-link' : 'btn-success')
    }


    setInterval(() => {
      router.replace(router.asPath, undefined, {scroll: false})
    }, 60_000 * 5)
  })

  return (
    <>
      <GithubCorner href="https://github.com/Jleagle/gym-tracker" bannerColor="#2f7ed8"/>
      <div className="row">
        <div className="col">

          <Gyms gyms={gyms}/>

          <h2>Now</h2>
          <LineChart data={now}/>

          <h2>By Hour</h2>
          <BarChart data={dayHour}/>
          <HeatMap data={weekHour}/>

          <h2>By Day</h2>
          <BarChart data={weekDay}/>
          <BarChart data={monthDay}/>

          <h2>By Month</h2>
          <BarChart data={yearMonth}/>

          <footer>
            Data updated every 10 minutes. If a gym has 10 or less members inside, it will show as 0.<br/>
            Created by <a href="https://jimeagle.com">Jim Eagle</a>.
          </footer>
        </div>
      </div>
    </>
  )
}

export default HomePage
