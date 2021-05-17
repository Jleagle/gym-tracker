import React from 'react'
import Link from 'next/link'
import {Button} from "react-bootstrap"

function Gyms(props) {

  const gyms = props.gyms.map((gym) => {
    return (
      <Button key={gym} type="button" variant="success">
        <Link href={'/' + gym}>{gym}</Link>
      </Button>
    )
  })

  gyms.push(
    <Button type="button" variant="link">
      <Link href="/new-gym">Add your gym!</Link>
    </Button>
  )

  return <div>{gyms}</div>
}

export default Gyms
