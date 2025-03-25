import { useState } from 'react';
import './InputForm.css';

interface InputFormProps {
  onStart: (name: string) => void;
}

function InputForm({ onStart }: InputFormProps) {
  const [inputText, setInputText] = useState('');
  const [submittedText, setSubmittedText] = useState<string | null>(null);
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
        // Call the onStart prop with the submitted name
        onStart(submittedText || inputText);
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
      
      {submittedText && !loading && (
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
    </div>
  );
}

export default InputForm; 