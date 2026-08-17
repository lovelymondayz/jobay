.PHONY: dev run pipeline stage psql clean

# Activate venv and run dev helper
dev:
	. .venv/bin/activate && python3 run_pipeline.py --stage scout

# Run full pipeline (all 7 stages)
pipeline:
	. .venv/bin/activate && python3 run_pipeline.py

# Run a single stage
# Usage: make stage STAGE=scout
stage:
	. .venv/bin/activate && python3 run_pipeline.py --stage $(STAGE)

# Docker
up:
	docker compose up -d postgres
	@echo "PostgreSQL running on port 5433"

down:
	docker compose down

# Database
psql:
	docker compose exec -it leadgen-postgres psql -U leadgen -d leadgen

# Individual stages
scout:
	. .venv/bin/activate && python3 run_pipeline.py --stage scout

evaluate:
	. .venv/bin/activate && python3 run_pipeline.py --stage evaluate

enrich:
	. .venv/bin/activate && python3 run_pipeline.py --stage enrich

build:
	. .venv/bin/activate && python3 run_pipeline.py --stage build

deploy:
	. .venv/bin/activate && python3 run_pipeline.py --stage deploy

outreach:
	. .venv/bin/activate && python3 run_pipeline.py --stage outreach

track:
	. .venv/bin/activate && python3 run_pipeline.py --stage track

# Utility
logs:
	docker compose logs -f

clean:
	docker compose down -v
	rm -rf __pycache__ *.pyc logs/*.log
