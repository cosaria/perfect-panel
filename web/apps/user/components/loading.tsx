interface LoadingProps {
  loading?: boolean;
}

export default function Loading({ loading = true }: LoadingProps) {
  return (
    <div
      className={`fixed inset-0 z-[9999] flex items-center justify-center transition-all duration-500 ${
        loading
          ? 'bg-background/95 pointer-events-auto opacity-100 backdrop-blur-[2px]'
          : 'pointer-events-none bg-transparent opacity-0 backdrop-blur-0'
      }`}
    >
      <div className='flex min-h-screen items-center justify-center'>
        <div className='text-center'>
          <div className='border-primary mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-b-2 border-t-2'></div>
          <p className='text-muted-foreground text-lg'>Loading...</p>
        </div>
      </div>
    </div>
  );
}
