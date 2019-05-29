import React, { Component, Fragment } from 'react';
import { reduxForm, getFormValues, SubmissionError } from 'redux-form';
import { connect } from 'react-redux';
import { get, map, isEmpty } from 'lodash';
import PropTypes from 'prop-types';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { replace } from 'react-router-redux';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import RadioButton from 'shared/RadioButton';
import './PPMPaymentRequest.css';
import Uploader from 'shared/Uploader';
import { createMoveDocument } from 'shared/Entities/modules/moveDocuments';
import Alert from 'shared/Alert';

export const missingWeightTickets = uploaders => {
  if (isEmpty(uploaders)) {
    return true;
  }
  const uploadersKeys = Object.keys(uploaders);
  for (const key of uploadersKeys) {
    // eslint-disable-next-line security/detect-object-injection
    if (uploaders[key].isEmpty()) {
      return true;
    }
  }
  return false;
};

class WeightTicket extends Component {
  state = {
    uploaderIsIdle: {},
    vehicleType: '',
    additionalWeightTickets: 'Yes',
  };
  uploaders = {};

  //  handleChange for vehicleType and additionalWeightTickets
  handleChange = (event, type) => {
    this.setState({ [type]: event.target.value });
  };

  onAddFile = uploaderName => () => {
    this.setState({
      uploaderIsIdle: { ...this.state.uploaderIsIdle, [uploaderName]: false },
    });
  };

  onChange = uploaderName => uploaderIsIdle => {
    this.setState({
      uploaderIsIdle: { ...this.state.uploaderIsIdle, [uploaderName]: uploaderIsIdle },
    });
  };

  submitWeightTickets = formValues => {
    const { moveId, currentPpm } = this.props;

    const moveDocumentSubmissions = [];
    const uploadersKeys = Object.keys(this.uploaders);
    for (const key of uploadersKeys) {
      // eslint-disable-next-line security/detect-object-injection
      let files = this.uploaders[key].getFiles();
      if (files.length > 0) {
        const uploadIds = map(files, 'id');
        const weightTicket = {
          moveId: moveId,
          personallyProcuredMoveId: currentPpm.id,
          uploadIds: uploadIds,
          title: key,
          moveDocumentType: 'WEIGHT_TICKET',
          notes: formValues.notes,
        };
        moveDocumentSubmissions.push(
          this.props.createMoveDocument(weightTicket).catch(() => {
            throw new SubmissionError({ _error: 'Error creating move document' });
          }),
        );
      }
    }
    return Promise.all(moveDocumentSubmissions)
      .then(() => {
        this.cleanup();
      })
      .catch(e => {
        throw new SubmissionError({ _error: e.message });
      });
  };

  cleanup = () => {
    const { reset } = this.props;
    const uploaders = this.uploaders;
    const uploadersKeys = Object.keys(this.uploaders);
    for (const key of uploadersKeys) {
      // eslint-disable-next-line security/detect-object-injection
      uploaders[key].clearFiles();
    }
    reset();
  };

  render() {
    const missingWeightTicket = missingWeightTickets(this.uploaders);
    const { additionalWeightTickets, vehicleType } = this.state;
    const { error, handleSubmit, submitting, schema } = this.props;
    const nextBtnLabel = additionalWeightTickets === 'Yes' ? 'Save & Add Another' : 'Save & Continue';
    const uploadEmptyTicketLabel =
      '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload upload empty weight ticket</span></span>';
    const uploadFullTicketLabel =
      '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload upload full weight ticket</span></span>';
    return (
      <Fragment>
        <WizardHeader
          title="Weight tickets"
          right={
            <ProgressTimeline>
              <ProgressTimelineStep name="Weight" current />
              <ProgressTimelineStep name="Expenses" />
              <ProgressTimelineStep name="Review" />
            </ProgressTimeline>
          }
        />
        <form>
          {error && (
            <div className="usa-grid">
              <div className="usa-width-one-whole error-message">
                <Alert type="error" heading="An error occurred">
                  Something went wrong contacting the server.
                </Alert>
              </div>
            </div>
          )}
          <div className="usa-grid">
            <SwaggerField
              fieldName="vehicle_options"
              swagger={schema}
              onChange={event => this.handleChange(event, 'vehicleType')}
              value={vehicleType}
              required
            />
            <SwaggerField fieldName="vehicle_nickname" swagger={schema} required />
            {(vehicleType === 'CAR' || vehicleType === 'BOX_TRUCK') && (
              <>
                <div className="dashed-divider" />

                <div className="usa-grid-full" style={{ marginTop: '1em' }}>
                  <div className="usa-width-one-third">
                    <strong className="input-header">Empty Weight</strong>
                    <SwaggerField
                      className="short-field"
                      fieldName="empty_weight"
                      swagger={schema}
                      hideLabel
                      required
                    />{' '}
                    lbs
                  </div>
                  <Uploader
                    options={{ labelIdle: uploadEmptyTicketLabel }}
                    onRef={ref => (this.uploaders['empty_weight'] = ref)}
                    onChange={this.onChange('empty_weight')}
                    onAddFile={this.onAddFile('empty_weight')}
                  />
                </div>

                <div className="usa-grid-full" style={{ marginTop: '1em' }}>
                  <div className="usa-width-one-third">
                    <strong className="input-header">Full Weight</strong>
                    <label className="full-weight-label">Full weight at destination</label>
                    <SwaggerField
                      className="short-field"
                      fieldName="full_weight"
                      swagger={schema}
                      hideLabel
                      required
                    />{' '}
                    lbs
                  </div>
                  <Uploader
                    options={{ labelIdle: uploadFullTicketLabel }}
                    onRef={ref => (this.uploaders['full_weight'] = ref)}
                    onChange={this.onChange('full_weight')}
                    onAddFile={this.onAddFile('full_weight')}
                  />
                </div>

                <SwaggerField fieldName="weight_ticket_date" swagger={schema} required />
                <div className="dashed-divider" />

                <div className="radio-group-wrapper">
                  <p className="radio-group-header">Do you have more weight tickets for another vehicle or trip?</p>
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="Yes"
                    value="Yes"
                    name="additional_weight_ticket"
                    checked={additionalWeightTickets === 'Yes'}
                    onChange={event => this.handleChange(event, 'additionalWeightTickets')}
                  />

                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="No"
                    value="No"
                    name="additional_weight_ticket"
                    shipment_line_items
                    checked={additionalWeightTickets === 'No'}
                    onChange={event => this.handleChange(event, 'additionalWeightTickets')}
                  />
                </div>
              </>
            )}

            {/* TODO: change onclick handler to go to next page in flow */}
            <PPMPaymentRequestActionBtns
              nextBtnLabel={nextBtnLabel}
              missingWeightTicket={missingWeightTicket}
              submitting={submitting}
              onClick={handleSubmit(this.submitWeightTickets)}
              displaySaveForLater={true}
            />
          </div>
        </form>
      </Fragment>
    );
  }
}

const formName = 'weight_ticket_wizard';
WeightTicket = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(WeightTicket);

WeightTicket.propTypes = {
  schema: PropTypes.object.isRequired,
};

function mapStateToProps(state, ownProps) {
  const moveId = ownProps.match.params.moveId;
  return {
    moveId: moveId,
    formValues: getFormValues(formName)(state),
    genericMoveDocSchema: get(state, 'swaggerInternal.spec.definitions.CreateGenericMoveDocumentPayload', {}),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    schema: get(state, 'swaggerInternal.spec.definitions.WeightTicketPayload', {}),
    currentPpm: get(state, 'ppm.currentPpm'),
  };
}
const mapDispatchToProps = {
  createMoveDocument,
  replace,
};

export default connect(mapStateToProps, mapDispatchToProps)(WeightTicket);
