import { useState } from 'react';
import './InputForm.css';

function InputForm() {
  const [inputText, setInputText] = useState('');
  const [submittedText, setSubmittedText] = useState<string | null>(null);
  const [started, setStarted] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (inputText.trim()) {
      setSubmittedText(inputText);
      // Here you could also send the data to an API
      console.log('Submitted:', inputText);
    }
  };

  const handleStart = () => {
    if (inputText.trim()) {
      setLoading(true);
      console.log('Loading started for:', inputText);
      
      // Simulate backend API call
      setTimeout(() => {
        setLoading(false);
        setStarted(true);
        console.log('Started with name:', inputText);
      }, 2000); // 2 seconds loading simulation
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
      
      {submittedText && !started && !loading && (
        <div className="result-container">
          <h3>Tên đã nhập:</h3>
          <p>{submittedText}</p>
          <button type="button" className="start-button" onClick={handleStart}>
            Bắt đầu
          </button>
        </div>
      )}

      {loading && (
        <div className="loading-container">
          <div className="loading-spinner"></div>
          <p>Đang xử lý...</p>
        </div>
      )}
      
      {started && (
        <div className="started-container">
          <h3>Đã bắt đầu!</h3>
          <p>Xin chào, {submittedText}!</p>
        </div>
      )}
    </div>
  );
}

export default InputForm; 