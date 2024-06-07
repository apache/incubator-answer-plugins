/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import React, { FC, cloneElement, useRef } from 'react';
import Spinner from 'react-bootstrap/Spinner';

interface EmbedContainerProps {
  children: React.ReactElement;
  height?: number | string;
}
const EmbedContainer: FC<EmbedContainerProps> = ({
  children,
  height = 350,
}) => {
  const loadingRef = useRef<HTMLSpanElement>(null);

  const handleLoad = () => {
    if (loadingRef.current) {
      const parentElement = loadingRef.current.parentElement;
      if (parentElement) {
        parentElement.style.height = height + 'px';
      }
      loadingRef.current.remove();
    }
  };
  let Component = children;
  if (children.type === 'iframe') {
    Component = cloneElement(children, { onLoad: handleLoad });
  }
  return (
    <>
      {Component}
      <span
        ref={loadingRef}
        className="loading position-absolute top-0 left-0 w-100 h-100 z-1 bg-white d-flex justify-content-center align-items-center">
        <Spinner animation="border" variant="secondary" />
      </span>
    </>
  );
};

export default EmbedContainer;
