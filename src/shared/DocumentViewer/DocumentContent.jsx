import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './index.css';

export class DocumentContent extends Component {
  imgEl = React.createRef();

  state = {
    imgHeight: 0,
  };

  adjustContainerHeight() {
    const imgHeight = this.imgEl.current.getBoundingClientRect().height;
    this.setState({
      imgHeight: imgHeight,
    });
  }

  render() {
    return (
      <div
        className="page"
        style={{
          minHeight: this.state.imgHeight + 50,
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
        }}
      >
        {this.props.contentType === 'application/pdf' ? (
          <div className="pdf-placeholder">
            {this.props.filename && <span className="filename">{this.props.filename}</span>}
            This PDF can be{' '}
            <a target="_blank" href={this.props.url}>
              viewed here
            </a>
            .
          </div>
        ) : (
          <div style={{ marginTop: this.state.imgHeight / 5, marginBottom: this.state.imgHeight / 5 }}>
            <img
              src={this.props.url}
              ref={this.imgEl}
              style={{ transform: `rotate(${this.props.orientation}deg)` }}
              onLoad={this.adjustContainerHeight.bind(this)}
              alt="document upload"
            />
          </div>
        )}

        <div>
          <button onClick={this.props.rotate.bind(this, this.props.uploadId, 'left')} data-direction="left">
            rotate left
          </button>
          <button onClick={this.props.rotate.bind(this, this.props.uploadId, 'right')} data-direction="right">
            rotate right
          </button>
        </div>
      </div>
    );
  }
}

DocumentContent.propTypes = {
  contentType: PropTypes.string.isRequired,
  filename: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
  uploadId: PropTypes.string.isRequired,
  orientation: PropTypes.string,
  rotate: PropTypes.func.isRequired,
};

export default DocumentContent;
