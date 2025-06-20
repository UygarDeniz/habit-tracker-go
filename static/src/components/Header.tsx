import { Link } from '@tanstack/react-router'
import { LogOut, Target } from 'lucide-react'
import { useUser } from '@/contexts/UserContext'
import { Button } from '@/components/ui/button'

export default function Header() {
  const { user, isAuthenticated, logout } = useUser()

  const handleLogout = async () => {
    await logout()
  }

  return (
    <>
      {/* Work in Progress Banner */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white text-center py-2 px-4">
        <p className="text-sm">
          🚧 <strong>Coming Soon!</strong> StreakCraft is currently in development. Sign up to be notified when we launch! 🚧
        </p>
      </div>

      {/* Header */}
      <header className="bg-white/80 backdrop-blur-sm border-b border-white/20 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
                <Target className="w-5 h-5 text-white" />
              </div>
              <Link to="/" className="text-xl font-bold text-gray-900 hover:text-gray-700 transition-colors">
                StreakCraft
              </Link>
            </div>
            <nav className="flex items-center space-x-4">
              {isAuthenticated ? (
                <div className="flex items-center space-x-4">
                  <div className="flex items-center space-x-2">
                    <img 
                      src={user?.picture} 
                      alt={user?.name} 
                      className="w-8 h-8 rounded-full"
                    />
                    <span className="text-gray-700 font-medium">{user?.name}</span>
                  </div>
                  <Button 
                    onClick={handleLogout} 
                    variant="outline" 
                    size="sm"
                    className="flex items-center space-x-1"
                  >
                    <LogOut className="w-4 h-4" />
                    <span>Logout</span>
                  </Button>
                </div>
              ) : (
                <Link 
                  to="/auth" 
                  className="text-gray-600 hover:text-gray-900 transition-colors font-medium"
                >
                  Login
                </Link>
              )}
            </nav>
          </div>
        </div>
      </header>
    </>
  )
} 