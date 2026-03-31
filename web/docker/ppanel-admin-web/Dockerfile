# Use the official lightweight Bun image
FROM oven/bun:latest AS base

# Set working directory
WORKDIR /app

# Create a non-root user for running the production application
RUN apt-get update \
 && apt-get install -y --no-install-recommends adduser \
 && rm -rf /var/lib/apt/lists/* \
 && addgroup --system --gid 1001 nodejs \
 && adduser --system --uid 1001 --ingroup nodejs --home /nonexistent --shell /usr/sbin/nologin nextjs

# Copy necessary files for production
COPY ./apps/admin/.next/standalone ./
COPY ./apps/admin/.next/static ./apps/admin/.next/static
COPY ./apps/admin/public ./apps/admin/public

# Change to non-root user
RUN chown -R nextjs:nodejs /app
USER nextjs

# Disable Next.js telemetry at runtime
ENV NEXT_TELEMETRY_DISABLED=1

# Set default command
CMD ["bun", "apps/admin/server.js"]
