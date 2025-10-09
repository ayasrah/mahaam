package mahaam.infra.monitor;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicBoolean;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.infra.Cache;
import mahaam.infra.Log;
import mahaam.infra.monitor.MonitorModel.Health;

public interface HealthService {
	void serverStarted(Health health);

	void startSendingPulses();

	void serverStopped();
}

@ApplicationScoped
class DefaultHealthService implements HealthService {

	@Inject
	HealthRepo healthRepo;

	private static final AtomicBoolean pulseSendingActive = new AtomicBoolean(false);
	private static CompletableFuture<Void> pulseTask;

	@Override
	public void serverStarted(Health health) {
		healthRepo.create(health);
	}

	@Override
	public void serverStopped() {
		Thread thread = new Thread(
				() -> {
					try {
						// Stop pulse sending
						pulseSendingActive.set(false);
						if (pulseTask != null) {
							pulseTask.cancel(true);
						}

						if (Cache.getHealthId() != null) {
							healthRepo.updateStopped(Cache.getHealthId());
						}
					} catch (Exception e) {
						Log.error(e.toString());
					}
				});
		thread.start();
	}

	@Override
	public void startSendingPulses() {
		if (pulseSendingActive.compareAndSet(false, true)) {
			pulseTask = CompletableFuture.runAsync(
					() -> {
						while (pulseSendingActive.get()) {
							try {
								if (Cache.getHealthId() != null) {
									healthRepo.updatePulse(Cache.getHealthId());
								}
								Thread.sleep(60000); // 1 minute
							} catch (InterruptedException e) {
								// Thread was interrupted, exit gracefully
								Thread.currentThread().interrupt();
								break;
							} catch (Exception e) {
								Log.error(e.toString());
							}
						}
					});
		}
	}
}
