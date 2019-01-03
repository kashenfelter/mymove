import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import { get } from 'lodash';
import moment from 'moment';

import Alert from 'shared/Alert';
import './InvoicePanel.css';

class InvoicePaymentAlert extends PureComponent {
  render() {
    let paymentAlert;
    const status = this.props.createInvoiceStatus;

    if (status.error) {
      //handle 409 status: shipment invoice already processed
      let httpResCode = get(status, 'error.response.status');
      let invoiceStatus = get(status, 'error.response.response.body.status');
      let aproverFirstName = get(status, 'error.response.response.body.approver_first_name');
      let aproverLastName = get(status, 'error.response.response.body.approver_last_name');
      let invoiceDate = moment(get(status, 'error.response.response.body.invoiced_date'));
      if (httpResCode === 409 && invoiceStatus === 'SUBMITTED') {
        paymentAlert = (
          <div>
            <Alert type="success" heading="Already approved">
              <span className="warning--header">
                {aproverFirstName} {aproverLastName} approved this invoice on {invoiceDate.format('DD-MMM-YYYY')} at{' '}
                {invoiceDate.format('kk:mm')}.
              </span>
            </Alert>
          </div>
        );
      } else if (httpResCode === 409 && (invoiceStatus === 'IN_PROCESS' || invoiceStatus === 'DRAFT')) {
        paymentAlert = (
          <div>
            <Alert type="success" heading="Already submitted">
              <span className="warning--header">
                {aproverFirstName} {aproverLastName} submitted this invoice on {invoiceDate.format('DD-MMM-YYYY')} at{' '}
                {invoiceDate.format('kk:mm')}.
              </span>
            </Alert>
          </div>
        );
      } else {
        paymentAlert = (
          <Alert type="error" heading="Oops, something went wrong!">
            <span className="warning--header">Please try again.</span>
          </Alert>
        );
      }
    } else if (status.isLoading) {
      paymentAlert = (
        <Alert type="loading" heading="Creating invoice">
          <span className="warning--header">Sending information to USBank/Syncada.</span>
        </Alert>
      );
    } else if (status.isSuccess) {
      paymentAlert = (
        <div>
          <Alert type="success" heading="Success!">
            <span className="warning--header">The invoice has been created and will be paid soon.</span>
          </Alert>
        </div>
      );
    }

    return <div>{paymentAlert}</div>;
  }
}

InvoicePaymentAlert.propTypes = {
  createInvoiceStatus: PropTypes.object,
};

export default InvoicePaymentAlert;
