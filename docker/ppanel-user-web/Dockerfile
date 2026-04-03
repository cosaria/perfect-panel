# Use the official lightweight Bun image
FROM oven/bun:latest AS base

# Set working directory
WORKDIR /app

# Create non-root user and set permissions
RUN apt-get update \
 && apt-get install -y --no-install-recommends adduser \
 && rm -rf /var/lib/apt/lists/* \
 && addgroup --system --gid 1001 nodejs \
 && adduser --system --uid 1001 --ingroup nodejs --home /nonexistent --shell /usr/sbin/nologin nextjs

# Copy build output and static files
COPY ./apps/user/.next/standalone ./
COPY ./apps/user/.next/static ./apps/user/.next/static
COPY ./apps/user/public ./apps/user/public

# Change to non-root user
RUN chown -R nextjs:nodejs /app
USER nextjs

# Disable Next.js telemetry
ENV NEXT_TELEMETRY_DISABLED=1

# Default command to start the server
CMD ["bun", "apps/user/server.js"]
