import React from 'react'
import {Alert, Button, FormLabel} from 'react-bootstrap'
import Link from 'next/link'
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome'
import {faSpinner} from "@fortawesome/free-solid-svg-icons"

function NewGymPage() {

  const submitForm = (event) => {

    event.preventDefault()

    const loading = document.querySelector('button[type=submit] svg')
    loading.classList.remove('d-none')

    // const url = 'http://localhost:' + process.env.PURE_PORT_BACKEND + '/new-gym'
    const url = 'https://gymtrackerapi.jimeagle.com/new-gym'

    const myHeaders = new Headers()
    myHeaders.append('Content-Type', 'application/json')

    fetch(url, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify([
        document.getElementById('email').value,
        document.getElementById('pin').value,
      ]),
    })
      .then(response => response.json())
      .then(response => {

        const alert = document.getElementById('alert')

        alert.classList.remove('d-none')

        if (response.success) {
          alert.classList.add('alert-success')
          alert.innerHTML = 'Success, thanks!'
        } else {
          alert.classList.add('alert-danger')
          alert.innerHTML = response.message
        }

        loading.classList.add('d-none')
      })

    return false
  }

  return (
    <>
      <p>Gym Tracker only works with PureGym.<br/>Please enter your email and PIN so we can access your gym's member count.</p>
      <Alert className="d-none" id="alert">.</Alert>

      <div className="row">
        <div className="col-sm-10 col-md-8 col-lg-6 col-xl-4">

          <form onSubmit={submitForm}>
            <div className="mb-3">
              <label htmlFor="email" className="form-label">Email Address</label>
              <input type="email" className="form-control" id="email"/>
            </div>

            <div className="mb-3">
              <FormLabel htmlFor="pin" className="form-label">PIN</FormLabel>
              <input type="number" pattern="[0-9]{6}" className="form-control" id="pin" maxLength="6" size="6"
                     required/>
            </div>

            <Button type="submit" className="btn btn-success">
              Submit <FontAwesomeIcon className="d-none" icon={faSpinner} spin/>
            </Button>
            <Link href="/"><a type="button" className="btn btn-primary float-end">Back</a></Link>
          </form>
        </div>
      </div>
    </>
  )
}

export default NewGymPage
