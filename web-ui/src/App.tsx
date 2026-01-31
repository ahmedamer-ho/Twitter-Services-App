import { AuthProvider, useAuth } from './context/AuthContext.tsx';
import { AuthFeature } from './features/AuthFeature.tsx';
import { MainTimeline } from './components/MainTimeline.tsx';

const AppContent: React.FC = () => {
  const { isAuthenticated } = useAuth();

  return (
    <div className="app-container">
      {isAuthenticated ? <MainTimeline /> : <AuthFeature />}
    </div>
  );
};

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;
