from flask import Flask, jsonify, request
import os

app = Flask(__name__)

@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({"status": "healthy", "service": "ml-services"})

@app.route('/predict', methods=['POST'])
def predict():
    data = request.json
    # Placeholder for prediction logic
    # In a real app, this would load a model and return predictions
    
    response = {
        "prediction": "positive",
        "confidence": 0.85,
        "input_processed": data
    }
    return jsonify(response)

@app.route('/anomaly', methods=['POST'])
def detect_anomaly():
    data = request.json
    # Placeholder for anomaly detection logic
    
    response = {
        "is_anomaly": False,
        "score": 0.12
    }
    return jsonify(response)

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 5000))
    app.run(host='0.0.0.0', port=port)
