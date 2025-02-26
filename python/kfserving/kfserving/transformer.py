# Copyright 2019 kubeflow.org.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import logging
from typing import List, Dict
from .server import Protocol
from .server import KFModel
from .server import KFSERVER_LOGLEVEL
from .server import PREDICTOR_URL_FORMAT
from kfserving.protocols.tensorflow_http import TensorflowRequestHandler
logging.basicConfig(level=KFSERVER_LOGLEVEL)


class Transformer(KFModel):
    def __init__(self, name: str,
                 predictor_host: str,
                 protocol: Protocol):
        super().__init__(name)
        self.predictor_url = PREDICTOR_URL_FORMAT.format(predictor_host, name)
        self.protocol = protocol
        self.ready = False

    def load(self):
        if not self.ready:
            self.ready = True

    # subclass of Transformer should implement preprocess
    def preprocess(self, inputs: Dict) -> Dict:
        return inputs

    # subclass of Transformer should implement postprocess
    def postprocess(self, inputs: List) -> List:
        return inputs

    def predict(self, inputs: List) -> List:
        if self.protocol == Protocol.tensorflow_http:
            return TensorflowRequestHandler.predict(inputs, self.predictor_url)
        else:
            raise NotImplementedError

    def explain(self, inputs: List) -> List:
        raise NotImplementedError

    def detect_outlier(self, inputs: List):
        raise NotImplementedError
