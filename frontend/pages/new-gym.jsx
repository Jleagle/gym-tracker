import React from 'react'
import {Button, FormLabel} from 'react-bootstrap'
import Link from 'next/link'

function NewGymPage() {

  const submitForm = () => {

    // const url = 'http://localhost:' + process.env.PURE_PORT_BACKEND + '/new-gym'
    const url = 'https://gymtrackerapi.jimeagle.com/new-gym'

    const myHeaders = new Headers()
    myHeaders.append('Content-Type', 'application/json')

    fetch(url, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
        email: document.getElementById('email').value,
        pin: document.getElementById('pin').value,
      }),
    }).then(response => response.json()).then(response => {
      console.log(response)
    })
  }

  return (
    <div className="row">
      <div className="col-sm-10 col-md-8 col-lg-6 col-xl-4">

        <p>xxx</p>

        <div className="mb-3">
          <label htmlFor="email" className="form-label">Email Address</label>
          <input type="email" className="form-control" id="email"/>
        </div>

        <div className="mb-3">
          <FormLabel htmlFor="pin" className="form-label">PIN</FormLabel>
          <input type="text" className="form-control" id="pin" maxLength="6" size="6"/>
        </div>

        <Button type="submit" className="btn btn-success" onClick={submitForm}>Submit</Button>
        <Link href="/"><a type="button" className="btn btn-primary float-end">Back</a></Link>

      </div>
    </div>
  )
}

export default NewGymPage
