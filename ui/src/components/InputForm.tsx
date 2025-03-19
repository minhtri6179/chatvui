import { useState } from 'react';
import './InputForm.css';

function InputForm() {
  const [inputText, setInputText] = useState('');
  const [submittedText, setSubmittedText] = useState<string | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (inputText.trim()) {
      setSubmittedText(inputText);
      // Here you could also send the data to an API
      console.log('Submitted:', inputText);
    }
  };

  return (
    <div className="input-form-container">
      <h2>Tên thân thương</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <input
            type="text"
            id="textInput"
            value={inputText}
            onChange={(e) => setInputText(e.target.value)}
            placeholder="Nhập tên của bạn..."
            className="text-input"
          />
        </div>
        <button type="submit" className="submit-button">
          Gửi
        </button>
      </form>
      
      {submittedText && (
        <div className="result-container">
          <h3>Tên đã nhập:</h3>
          <p>{submittedText}</p>
        </div>
      )}
    </div>
  );
}

export default InputForm; 