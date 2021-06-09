import React from 'react'
import Link from 'next/link'
import {Button} from "react-bootstrap"

function Gyms(props) {

  const gyms = props.gyms.map((gym) => {
    return (
      <Link key={gym} href={'/' + gym}>
        <Button type="button" variant="link" className="me-2" data-gym={gym}>
          {gym}
        </Button>
      </Link>
    )
  })

  gyms.push(
    <Link key="new" href="/new-gym">
      <Button type="button" variant="link" className="me-2">
        Add your gym!
      </Button>
    </Link>
  )

  gyms.unshift(
    <Link key="all" href="/all">
      <Button type="button" variant="link" className="me-2" data-gym="all">
        all
      </Button>
    </Link>
  )

  return <div id="gyms">{gyms}</div>
}

export default Gyms
